package router

import v1 "go-admin/api/v1"

type routerGroup struct {
	CaptchaRouter
	DBRouter
}

var RouterGroup = new(routerGroup)

var (
	captchaApi = v1.ApiGroup.CaptchaApi
	dbApi      = v1.ApiGroup.DBApi
)
