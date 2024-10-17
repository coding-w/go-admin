package router

import v1 "go-admin/api/v1"

type routerGroup struct {
	CaptchaRouter
	DBRouter
	ApiRouter
	AuthorityRouter
	AuthorityBtnRouter
	SysDictionaryRouter
	SysMenuRouter
	SysUserRouter
}

var RouterGroup = new(routerGroup)

var (
	captchaApi      = v1.ApiGroup.CaptchaApi
	dbApi           = v1.ApiGroup.DBApi
	userApi         = v1.ApiGroup.UserApi
	menuApi         = v1.ApiGroup.AuthMenuApi
	authApi         = v1.ApiGroup.AuthorityApi
	apiApi          = v1.ApiGroup.SysApiApi
	authorityBtnApi = v1.ApiGroup.AuthorityBtnApi
	dictionaryApi   = v1.ApiGroup.DictionaryApi
)
