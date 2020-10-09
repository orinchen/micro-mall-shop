package main

import (
	"gitee.com/cristiane/micro-mall-shop/vars"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/cristiane/micro-mall-shop/startup"
	"gitee.com/kelvins-io/kelvins/app"
)

const APP_NAME = "micro-mall-shop"

func main() {
	application := &kelvins.GRPCApplication{
		Application: &kelvins.Application{
			LoadConfig: startup.LoadConfig,
			SetupVars:  startup.SetupVars,
			Name:       APP_NAME,
		},
		RegisterGRPCServer: startup.RegisterGRPCServer,
		RegisterGateway:    startup.RegisterGateway,
		RegisterHttpRoute:  startup.RegisterHttpRoute,
	}
	app.RunGRPCApplication(application)
	vars.App = application
}