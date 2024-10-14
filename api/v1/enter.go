package v1

import "go-admin/service"

type apiGroup struct {
	CaptchaApi
	DBApi
	//UserApi
	//AuthMenuApi
	//AuthorityApi
	//SysApiApi
	//AuthorityBtnApi
	//AutoCodeApi
	//DictionaryApi
}

var ApiGroup = new(apiGroup)

var (
	initDBService = service.ServiceGroup.InitDBService
)
