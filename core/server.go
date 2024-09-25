package core

import (
	"fmt"
	"go-admin/global"
	"go-admin/initialize"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunServer() {
	if global.GA_CONFIG.System.UseMultipoint || global.GA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}

	// 初始化路由
	Router := initialize.Routers()
	address := fmt.Sprintf(":%d", global.GA_CONFIG.System.Port)
	s := initServer(address, Router)

	global.GA_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`欢迎使用 go-admin`, address)
	global.GA_LOG.Error(s.ListenAndServe().Error())
}
