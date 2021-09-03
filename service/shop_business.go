package service

import (
	"context"
	"fmt"
	"gitee.com/cristiane/micro-mall-shop/model/args"
	"gitee.com/cristiane/micro-mall-shop/model/mysql"
	"gitee.com/cristiane/micro-mall-shop/pkg/code"
	"gitee.com/cristiane/micro-mall-shop/pkg/util"
	"gitee.com/cristiane/micro-mall-shop/pkg/util/email"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_pay_proto/pay_business"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_search_proto/search_business"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_shop_proto/shop_business"
	"gitee.com/cristiane/micro-mall-shop/repository"
	"gitee.com/cristiane/micro-mall-shop/vars"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

func CreateShopBusiness(ctx context.Context, req *shop_business.ShopApplyRequest) (shopId int64, retCode int) {
	retCode = code.Success
	if req.OperationType == shop_business.OperationType_CREATE {
		// 创建账户前检查是为了防止数据表没有唯一性检查
		exist, err := repository.CheckShopBusinessInfoExist(int(req.MerchantId), req.NickName)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CheckShopBusinessInfoExistByMerchantId err: %v, MerchantId: %v", err, req.MerchantId)
			retCode = code.ErrorServer
			return
		}
		if exist {
			retCode = code.ShopBusinessExist
			return
		}
		shopCode := uuid.New().String()
		shopInfo := mysql.ShopBusinessInfo{
			NickName:         req.NickName,
			FullName:         req.FullName,
			ShopCode:         shopCode,
			RegisterAddr:     req.RegisterAddr,
			BusinessAddr:     req.BusinessAddr,
			LegalPerson:      req.MerchantId,
			BusinessLicense:  req.BusinessLicense,
			TaxCardNo:        req.TaxCardNo,
			BusinessDesc:     req.BusinessDesc,
			SocialCreditCode: req.SocialCreditCode,
			OrganizationCode: req.OrganizationCode,
			State:            2, // 线上流程禁止
			CreateTime:       time.Now(),
			UpdateTime:       time.Now(),
		}
		tx := kelvins.XORM_DBEngine.NewSession()
		err = tx.Begin()
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo NewSession err: %v", err)
			retCode = code.ErrorServer
			return
		}
		// 创建店铺账户
		err = repository.CreateShopBusinessInfo(tx, &shopInfo)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo Rollback err: %v", errRollback)
			}
			if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
				retCode = code.ShopBusinessExist
				return
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo err: %v, shopInfo: %v", err, json.MarshalToStringNoError(shopInfo))
			retCode = code.ErrorServer
			return
		}
		serverName := args.RpcServiceMicroMallPay
		conn, err := util.GetGrpcClient(ctx,serverName)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo Rollback err: %v", errRollback)
			}
			kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
			retCode = code.ErrorServer
			return
		}
		//defer conn.Close()
		client := pay_business.NewPayBusinessServiceClient(conn)
		balance := "1.9999"
		accountReq := pay_business.CreateAccountRequest{
			Owner:       shopCode,
			AccountType: pay_business.AccountType_Company,
			CoinType:    pay_business.CoinType_CNY,
			Balance:     balance,
		}
		rsp, err := client.CreateAccount(ctx, &accountReq)
		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo Rollback err: %v", errRollback)
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount %v,err: %v", serverName, err)
			retCode = code.ErrorServer
			return
		}
		if rsp.Common.Code != pay_business.RetCode_SUCCESS {
			errRollback := tx.Rollback()
			if errRollback != nil {
				kelvins.ErrLogger.Errorf(ctx, "CreateShopBusinessInfo Rollback err: %v", errRollback)
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateAccount req %v,rsp: %v", json.MarshalToStringNoError(req), json.MarshalToStringNoError(rsp))
			retCode = code.ErrorServer
			return
		}
		err = tx.Commit()
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "Commit err: %v", err)
			retCode = code.ErrorServer
			return
		}
		shopIdInfo, err := repository.GetShopBusinessInfo("shop_id", shopCode)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetShopBusinessInfo err: %v, shopCode: %v", err, shopCode)
			retCode = code.ErrorServer
			return
		}
		shopId = shopIdInfo.ShopId
		kelvins.GPool.SendJob(func() {
			// 发送登录邮件
			emailNotice := fmt.Sprintf(args.UserApplyShopTemplate, req.MerchantId, time.Now().String(), req.FullName, balance)
			for _, receiver := range vars.EmailNoticeSetting.Receivers {
				err = email.SendEmailNotice(ctx, receiver, vars.App.Name, emailNotice)
				if err != nil {
					kelvins.ErrLogger.Info(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
					return
				}
			}
		})
	} else if req.OperationType == shop_business.OperationType_UPDATE {
		shopId = req.ShopId
		query := map[string]interface{}{}
		if req.ShopId > 0 {
			query["shop_id"] = req.ShopId
		}
		if req.MerchantId > 0 {
			query["legal_person"] = req.MerchantId
		}
		maps := map[string]interface{}{
			"nick_name":     req.NickName,
			"full_name":     req.FullName,
			"register_addr": req.RegisterAddr,
			"business_addr": req.BusinessAddr,
			"business_desc": req.BusinessDesc,
		}
		err := repository.UpdateShopBusinessInfo(query, maps)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "UpdateShopBusinessInfo err: %v, query: %v, maps: %v", err, json.MarshalToStringNoError(query), json.MarshalToStringNoError(maps))
			retCode = code.ErrorServer
			return
		}

		kelvins.GPool.SendJob(func() {
			// 发送登录邮件
			emailNotice := fmt.Sprintf(args.UserModifyShopTemplate, req.MerchantId, time.Now().String(), req.FullName)
			for _, receiver := range vars.EmailNoticeSetting.Receivers {
				err = email.SendEmailNotice(ctx, receiver, vars.App.Name, emailNotice)
				if err != nil {
					kelvins.ErrLogger.Info(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
					return
				}
			}
		})
	} else if req.OperationType == shop_business.OperationType_DELETE {
		shopId = req.ShopId
		where := map[string]interface{}{
			"shop_id": req.ShopId,
		}
		if req.MerchantId > 0 {
			where["legal_person"] = req.MerchantId
		}
		err := repository.DeleteShopBusinessInfo(where)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "DeleteShopBusinessInfo err: %v, where: %v", err, json.MarshalToStringNoError(where))
			retCode = code.ErrorServer
			return
		}

		kelvins.GPool.SendJob(func() {
			// 发送登录邮件
			emailNotice := fmt.Sprintf(args.UserCloseShopTemplate, req.MerchantId, time.Now().String(), req.ShopId)
			for _, receiver := range vars.EmailNoticeSetting.Receivers {
				err = email.SendEmailNotice(ctx, receiver, vars.App.Name, emailNotice)
				if err != nil {
					kelvins.ErrLogger.Info(ctx, "SendEmailNotice err %v, emailNotice: %v", err, emailNotice)
					return
				}
			}
		})
		return
	}
	return
}

func GetShopMaterial(ctx context.Context, shopId int64) (*mysql.ShopBusinessInfo, int) {
	shopInfo, err := repository.GetShopBusinessInfoByShopId(shopId)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetShopBusinessInfoByShopId err: %v, shopId: %v", err, shopId)
		return shopInfo, code.ErrorServer
	}
	return shopInfo, code.Success
}

func SearchShopSync(ctx context.Context, shopId int64, pageSize, pageNum int) ([]*shop_business.SearchSyncShopEntry, int) {
	result := make([]*shop_business.SearchSyncShopEntry, 0)
	var shopIds []int64
	if shopId > 0 {
		shopIds = append(shopIds, shopId)
	}
	shopInfoList, err := repository.GetShopInfoList("*", shopIds, pageSize, pageNum)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "SearchShopSync err: %v, shopId: %v", err, shopId)
		return result, code.ErrorServer
	}
	result = make([]*shop_business.SearchSyncShopEntry, len(shopInfoList))
	for i := 0; i < len(shopInfoList); i++ {
		entry := &shop_business.SearchSyncShopEntry{
			ShopId:       shopInfoList[i].ShopId,
			NickName:     shopInfoList[i].NickName,
			FullName:     shopInfoList[i].FullName,
			ShopCode:     shopInfoList[i].ShopCode,
			RegisterAddr: shopInfoList[i].RegisterAddr,
			BusinessAddr: shopInfoList[i].BusinessAddr,
			BusinessDesc: shopInfoList[i].BusinessDesc,
		}
		result[i] = entry
	}
	return result, code.Success
}

func GetShopInfoList(ctx context.Context, shopIds []int64) ([]mysql.ShopBusinessInfo, int) {
	shopInfoList, err := repository.GetShopInfoList("*", shopIds, 0, 0)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetShopInfo err: %v, shopIds: %v", err, json.MarshalToStringNoError(shopIds))
		return shopInfoList, code.ErrorServer
	}
	return shopInfoList, code.Success
}

func SearchShop(ctx context.Context, req *shop_business.SearchShopRequest) (result []*shop_business.SearchShopInfo, retCode int) {
	result = make([]*shop_business.SearchShopInfo, 0)
	retCode = code.Success
	serverName := args.RpcServiceMicroMallSearch
	conn, err := util.GetGrpcClient(ctx, serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		retCode = code.ErrorServer
		return
	}
	//defer conn.Close()
	client := search_business.NewSearchBusinessServiceClient(conn)
	searchReq := search_business.ShopSearchRequest{
		ShopKey: req.Keyword,
	}
	searchRsp, err := client.ShopSearch(ctx, &searchReq)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "ShopSearch %v,err: %v", serverName, err)
		retCode = code.ErrorServer
		return
	}
	if searchRsp.Common.Code != search_business.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "ShopSearch  req: %+v, rsp: %+v", json.MarshalToStringNoError(searchReq), json.MarshalToStringNoError(searchRsp))
		retCode = code.ErrorServer
		return
	}
	if len(searchRsp.List) == 0 {
		return
	}
	shopIds := make([]int64, 0)
	for i := range searchRsp.List {
		if searchRsp.List[i].ShopId == "" {
			continue
		}
		shopId, err := strconv.ParseInt(searchRsp.List[i].ShopId, 10, 64)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "ShopSearch  ParseInt err: %v, shopId: %s", err, searchRsp.List[i].ShopId)
			retCode = code.ErrorServer
		}
		if shopId > 0 {
			shopIds = append(shopIds, shopId)
		}
	}
	shopList, retCode := GetShopInfoList(ctx, shopIds)
	if retCode != code.Success {
		return
	}
	if len(shopList) == 0 {
		return
	}
	shopIdToInfo := map[int64]mysql.ShopBusinessInfo{}
	for i := 0; i < len(shopList); i++ {
		shopIdToInfo[shopList[i].ShopId] = shopList[i]
	}
	result = make([]*shop_business.SearchShopInfo, 0, len(searchRsp.List))
	for i := range searchRsp.List {
		if searchRsp.List[i].ShopId == "" {
			continue
		}
		shopId, err := strconv.ParseInt(searchRsp.List[i].ShopId, 10, 64)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "SearchShop ParseInt err: %v, shopId: %v", err, searchRsp.List[i].ShopId)
			retCode = code.ErrorServer
			return
		}
		if _, ok := shopIdToInfo[shopId]; !ok {
			continue
		}
		entry := &shop_business.SearchShopInfo{
			Info: &shop_business.ShopMaterial{
				ShopId:           shopIdToInfo[shopId].ShopId,
				MerchantId:       shopIdToInfo[shopId].LegalPerson,
				NickName:         shopIdToInfo[shopId].NickName,
				FullName:         shopIdToInfo[shopId].FullName,
				RegisterAddr:     shopIdToInfo[shopId].RegisterAddr,
				BusinessAddr:     shopIdToInfo[shopId].BusinessAddr,
				BusinessLicense:  shopIdToInfo[shopId].BusinessLicense,
				TaxCardNo:        shopIdToInfo[shopId].TaxCardNo,
				BusinessDesc:     shopIdToInfo[shopId].BusinessDesc,
				SocialCreditCode: shopIdToInfo[shopId].SocialCreditCode,
				OrganizationCode: shopIdToInfo[shopId].OrganizationCode,
				ShopCode:         shopIdToInfo[shopId].ShopCode,
			},
			Score: searchRsp.List[i].Score,
		}
		result = append(result, entry)
	}

	return result, retCode
}

const sqlSelectShopMajorInfo = "shop_id,nick_name,shop_code,state"

func GetShopMajorInfo(ctx context.Context, req *shop_business.GetShopMajorInfoRequest) (result []*shop_business.ShopMajorInfo, retCode int) {
	retCode = code.Success
	shopInfoList, err := repository.GetShopInfoList(sqlSelectShopMajorInfo, req.ShopIds, 0, 0)
	result = make([]*shop_business.ShopMajorInfo, len(shopInfoList))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetShopInfoList err: %v, shopId: %s", err, req.GetShopIds())
		retCode = code.ErrorServer
		return
	}
	if len(shopInfoList) == 0 {
		retCode = code.ShopBusinessNotExist
		return
	}
	if len(shopInfoList) != len(req.GetShopIds()) {
		retCode = code.ShopBusinessNotExist
		return
	}
	for i := 0; i < len(shopInfoList); i++ {
		if shopInfoList[i].State != 2 {
			retCode = code.ShopBusinessStateNotVerify
			return
		}
		majorInfo := &shop_business.ShopMajorInfo{
			ShopId:   shopInfoList[i].ShopId,
			ShopCode: shopInfoList[i].ShopCode,
			ShopName: shopInfoList[i].NickName,
		}
		result[i] = majorInfo
	}
	return
}
