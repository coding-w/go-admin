package service

import "go-admin/service/initdb"

type serviceGroup struct {
	JwtService
	initdb.InitDBService
	ApiService
	CasbinService
	AuthorityBtnService
	AuthorityService
	MenuService
	BaseMenuService
	DictionaryService
	UserService
}

var ServiceGroup = new(serviceGroup)
