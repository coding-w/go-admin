package initialize

import (
	"github.com/patrickmn/go-cache"
	"go-admin/global"
	"go-admin/utils"
	"time"
)

func OtherInit() {
	dr, err := utils.ParseDuration(global.GA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}
	_, err = utils.ParseDuration(global.GA_CONFIG.JWT.BufferTime)
	if err != nil {
		panic(err)
	}
	cleanupInterval := dr + time.Hour
	// 初始化 缓存
	global.LocalCache = cache.New(dr, cleanupInterval)
}
