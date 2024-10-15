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
	OperationRecordService
}

var ServiceGroup = new(serviceGroup)
