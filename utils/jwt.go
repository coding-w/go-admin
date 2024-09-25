package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go-admin/global"
	"go-admin/model/dto"
	"time"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("couldn't handle this token")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GA_CONFIG.JWT.SigningKey),
	}
}

// CreateClaims
// @description  : 创建一个claims
func (j *JWT) CreateClaims(baseClaims dto.BaseClaims) dto.CustomClaims {
	bf, _ := ParseDuration(global.GA_CONFIG.JWT.BufferTime)
	ep, _ := ParseDuration(global.GA_CONFIG.JWT.ExpiresTime)
	return dto.CustomClaims{
		BaseClaims: baseClaims,
		BufferTime: int64(bf / time.Second), // 缓冲时间1天 缓冲时间内会获得新的token刷新令牌 此时一个用户会存在两个有效令牌 但是前端只留一个 另一个会丢失
		RegisteredClaims: jwt.RegisteredClaims{
			Audience:  jwt.ClaimStrings{"xa"},                    // 受众
			NotBefore: jwt.NewNumericDate(time.Now().Add(-1000)), // 签名生效时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ep)),    // 过期时间 7天  配置文件
			Issuer:    global.GA_CONFIG.JWT.Issuer,               // 签名的发行者
		},
	}
}

// CreateToken 生成JWT Token
func (j *JWT) CreateToken(claims dto.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// CreateTokenByOldToken 旧token 换新token 使用归并回源避免并发问题
func (j *JWT) CreateTokenByOldToken(oldToken string, claims dto.CustomClaims) (string, error) {
	// 使用 singleflight.Group 的 Do 方法，确保同一时间只会执行一次给定 key 的函数调用
	v, err, _ := global.GA_Concurrency_Control.Do("JWT:"+oldToken, func() (interface{}, error) {
		return j.CreateToken(claims)
	})
	// 将结果转换为 string 并返回，同时返回可能的错误
	return v.(string), err
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*dto.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &dto.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		// token 不是 nil，进一步检查 token.Claims 是否是 *dto.CustomClaims 类型，并且 token 是否有效
		if claims, ok := token.Claims.(*dto.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid
	}
}
