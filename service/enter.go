package service

type serviceGroup struct {
	JwtService
}

var ServiceGroup = new(serviceGroup)
