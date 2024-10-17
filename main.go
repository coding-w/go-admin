package main

import (
	"go-admin/core"
	"go-admin/global"
	"go-admin/initialize"
	"go.uber.org/zap"
)

func main() {
	// 读取配置
	global.GA_VIPER = core.Viper()
	// 初始化 日志系统
	global.GA_LOG = core.Zap()
	zap.ReplaceGlobals(global.GA_LOG)
	// 初始化 gorm
	global.GA_DB = initialize.GormInit()
	initialize.OtherInit()
	if global.GA_DB != nil {
		// 初始化表
		initialize.RegisterTables()
		// 程序结束前关闭数据库链接
		db, _ := global.GA_DB.DB()
		defer db.Close()
	}
	// 启动 定时任务
	initialize.Timer()
	// 启动服务
	core.RunServer()
}
