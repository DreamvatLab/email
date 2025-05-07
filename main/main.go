package main

import (
	"github.com/DeamvatLab/email"
	"github.com/DeamvatLab/email/main/core"
	"github.com/DeamvatLab/email/main/svc"
	"github.com/DreamvatLab/go/xerr"
	"github.com/DreamvatLab/host/hconsul"
)

func main() {
	// 登记服务信息
	hconsul.RegisterServiceInfo(core.ConfigProvider)

	// 注册服务
	email.RegisterEmailServiceServer(core.Host.GetGRPCServer(), &svc.EmailService{})

	// 运行
	xerr.FatalIfErr(core.Host.Run())
}
