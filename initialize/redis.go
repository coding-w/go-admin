package initialize

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go-admin/global"
	"go.uber.org/zap"
)

func Redis() {
	redisConfig := global.GA_CONFIG.Redis

	var client redis.UniversalClient

	client = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	ping, err := client.Ping(context.Background()).Result()
	if err != nil {
		global.GA_LOG.Error("redis connect ping failed, err:", zap.Error(err))
		panic(err)
	} else {
		global.GA_LOG.Info("redis connect ping result:", zap.String("ping", ping))
		global.GA_REDIS = client
	}

}
