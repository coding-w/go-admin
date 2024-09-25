package router

import "github.com/gin-gonic/gin"

type CaptchaRouter struct {
}

func (br *CaptchaRouter) InitCaptchaRouter(Router *gin.RouterGroup) {
	Router.POST("captcha", captchaApi.Captcha)
}
