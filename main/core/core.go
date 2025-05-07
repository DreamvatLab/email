package core

import (
	"github.com/DreamvatLab/go/xconfig"
	"github.com/DreamvatLab/go/xerr"
	"github.com/DreamvatLab/go/xlog"
	"github.com/DreamvatLab/host/hconsul"
	"github.com/DreamvatLab/host/hgrpc"
	"github.com/DreamvatLab/host/hlog"
)

var (
	Host           hgrpc.IGRPCServiceHost
	ConfigProvider xconfig.IConfigProvider
)

func init() {
	cp := xconfig.NewJsonConfigProvider()
	Host = hgrpc.NewGRPCServiceHost(cp)
	ConfigProvider = Host.GetConfigProvider()

	// 获取日志配置
	var logConfig *xlog.LogConfig
	err := ConfigProvider.GetStruct("Log", &logConfig)
	xerr.FatalIfErr(err)

	// 获取consul配置
	var consulConfig *hconsul.ConsulConfig
	err = ConfigProvider.GetStruct("Consul", &consulConfig)
	xerr.FatalIfErr(err)

	// 获取日志客户端ID
	logClientID := ConfigProvider.GetString("Log.ClientID")

	// 创建grpc sink
	grpcSink := hlog.NewGrpcSink(consulConfig, logClientID)

	// 初始化日志, 添加GRPC远程sink
	xlog.Init(logConfig, grpcSink)

}
