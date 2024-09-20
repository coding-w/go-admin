package constants

import "time"

// EXPIRATION 缓存有效期，默认720（分钟）
const EXPIRATION = 720 * time.Minute

// CAPTCHA_EXPIRATION 验证码有效期，默认5分钟
const CAPTCHA_EXPIRATION = 5 * time.Minute
