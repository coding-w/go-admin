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
}
