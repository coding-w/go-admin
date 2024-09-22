package router

import "go-admin/router/system"

type routerGroup struct {
	System system.RouterGroup
}

var RouterGroupApp = new(routerGroup)
