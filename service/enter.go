package service

import "go-admin/service/initdb"

type serviceGroup struct {
	JwtService
	initdb.InitDBService
}

var ServiceGroup = new(serviceGroup)
