package main

import (
	"go-admin/core"
	"go-admin/global"
	"go-admin/initialize"
	"go.uber.org/zap"
)

func main() {
	global.GA_VIPER = core.Viper()
	global.GA_LOG = core.Zap()
	zap.ReplaceGlobals(global.GA_LOG)
	global.GA_DB = initialize.GormInit()
	initialize.OtherInit()
	if global.GA_DB != nil {
		// 程序结束前关闭数据库链接
		db, _ := global.GA_DB.DB()
		defer db.Close()
	}
	// 启动服务
	core.RunServer()
}
