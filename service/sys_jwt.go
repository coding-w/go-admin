package service

import (
	"context"
	"go-admin/global"
	"go-admin/model/system"
	"go-admin/utils"
)

type JwtService struct{}

func (js *JwtService) IsBlacklist(token string) bool {
	_, ok := global.LocalCache.Get(token)
	return ok
}

// GetRedisJWT 从redis获取jwt
func (js *JwtService) GetRedisJWT(username string) (redisJWT string, err error) {
	redisJWT, err = global.GA_REDIS.Get(context.Background(), username).Result()
	return redisJWT, err
}

// SetRedisJWT 设置jwt缓存 并设置过期时间
func (js *JwtService) SetRedisJWT(token string, username string) error {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(global.GA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return err
	}
	timer := dr
	return global.GA_REDIS.Set(context.Background(), username, token, timer).Err()
}

func (js *JwtService) JsonInBlacklist(jwt system.JwtBlacklist) (err error) {
	err = global.GA_DB.Create(&jwt).Error
	if err != nil {
		return err
	}
	global.LocalCache.SetDefault(jwt.Jwt, struct{}{})
	return
}
