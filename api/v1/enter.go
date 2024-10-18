package v1

import "go-admin/service"

type apiGroup struct {
	CaptchaApi
	DBApi
	UserApi
	AuthMenuApi
	AuthorityApi
	SysApiApi
	AuthorityBtnApi
	DictionaryApi
}

var ApiGroup = new(apiGroup)

var (
	initDBService       = service.ServiceGroup.InitDBService
	apiService          = service.ServiceGroup.ApiService
	casbinService       = service.ServiceGroup.CasbinService
	authorityBtnService = service.ServiceGroup.AuthorityBtnService
	authorityService    = service.ServiceGroup.AuthorityService
	menuService         = service.ServiceGroup.MenuService
	baseMenuService     = service.ServiceGroup.BaseMenuService
	dictionaryService   = service.ServiceGroup.DictionaryService
	userService         = service.ServiceGroup.UserService
	jwtService          = service.ServiceGroup.JwtService
	fileUploadService   = service.ServiceGroup.FileUploadService
)
