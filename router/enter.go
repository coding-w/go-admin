package router

import v1 "go-admin/api/v1"

type routerGroup struct {
	CaptchaRouter
}

var RouterGroup = new(routerGroup)

var (
	captchaApi = v1.ApiGroup.CaptchaApi
)
