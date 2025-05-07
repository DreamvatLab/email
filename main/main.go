package main

import (
	"github.com/DreamvatLab/email"
	"github.com/DreamvatLab/email/main/core"
	"github.com/DreamvatLab/email/main/svc"
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
