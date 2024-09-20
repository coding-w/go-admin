package global

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go-admin/config"
	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var (
	GA_CONFIG              *config.Config // 全局配置
	GA_VIPER               *viper.Viper
	GA_LOG                 *zap.Logger
	GA_DB                  *gorm.DB
	GA_REDIS               redis.UniversalClient
	GA_ROUTERS             gin.RoutesInfo
	LocalCache             *cache.Cache
	GA_Concurrency_Control = &singleflight.Group{}
)
