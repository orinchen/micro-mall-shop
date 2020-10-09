package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-shop/model/mysql"
	"gitee.com/cristiane/micro-mall-shop/pkg/code"
	"gitee.com/cristiane/micro-mall-shop/proto/micro_mall_users_proto/users"
	"gitee.com/cristiane/micro-mall-shop/repository"
	"gitee.com/kelvins-io/common/errcode"
	"gitee.com/kelvins-io/kelvins"
	"strings"
	"time"
)

func MerchantsMaterial(ctx context.Context, req *users.MerchantsMaterialRequest) (int, int) {
	var merchantId int
	exist, err := repository.CheckUserExistById(int(req.Info.Uid))
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CheckUserExistById err: %v,req : %+v", err, req)
		return merchantId, code.ErrorServer
	}
	if !exist {
		return merchantId, code.UserNotExist
	}
	if req.OperationType == users.OperationType_CREATE {
		merchantMaterial := mysql.MerchantInfo{
			Uid:          int(req.Info.Uid),
			RegisterAddr: req.Info.RegisterAddr,
			HealthCardNo: req.Info.HealthCardNo,
			Identity:     int(req.Info.Identity),
			State:        int(req.Info.State),
			TaxCardNo:    req.Info.TaxCardNo,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err := repository.CreateMerchantsMaterial(&merchantMaterial)
		if err != nil {
			if strings.Contains(err.Error(), errcode.GetErrMsg(code.DBDuplicateEntry)) {
				return merchantId, code.MerchantExist
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateMerchantsMaterial err: %v,merchantMaterial : %+v", err, merchantMaterial)
			return merchantId, code.ErrorServer
		}
		record, err := repository.GetMerchantsMaterialByUid(int(req.Info.Uid))
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "GetMerchantsMaterialByUid err: %v,uid : %+v", err, req.Info.Uid)
			return merchantId, code.ErrorServer
		}
		merchantId = record.MerchantId
		return merchantId, code.Success
	} else if req.OperationType == users.OperationType_UPDATE {
		query := map[string]interface{}{
			"uid": req.Info.Uid,
		}
		maps := map[string]interface{}{
			"register_addr":  req.Info.RegisterAddr,
			"health_card_no": req.Info.HealthCardNo,
			"identity":       req.Info.Identity,
			"state":          req.Info.State,
			"tax_card_no":    req.Info.TaxCardNo,
			"update_time":    time.Now(),
		}
		err := repository.UpdateMerchantsMaterial(query, maps)
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "UpdateMerchantsMaterial err: %v,query : %+v, maps: %+v", err, query, maps)
			return merchantId, code.ErrorServer
		}
		return merchantId, code.Success
	}
	return merchantId, code.Success
}