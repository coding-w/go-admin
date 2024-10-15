package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/middleware"
)

type SysMenuRouter struct {
}

func (smr *SysMenuRouter) InitSysMenuRouter(PrivateRouter *gin.RouterGroup) {

	menuRouter := PrivateRouter.Group("menu").Use(middleware.OperationRecord())
	menuRouterWithoutRecord := PrivateRouter.Group("menu")
	{
		menuRouter.POST("addBaseMenu", menuApi.AddMenu)               // 新增菜单
		menuRouter.POST("addMenuAuthority", menuApi.AddMenuAuthority) // 增加menu和角色关联关系
		menuRouter.POST("deleteBaseMenu", menuApi.DeleteMenu)         // 删除菜单
		menuRouter.POST("updateBaseMenu", menuApi.UpdateMenu)         // 更新菜单
	}
	{
		menuRouterWithoutRecord.POST("getMenu", menuApi.GetMenu)                   // 获取菜单树
		menuRouterWithoutRecord.POST("getMenuList", menuApi.GetMenuList)           // 分页获取基础menu列表
		menuRouterWithoutRecord.POST("getBaseMenuTree", menuApi.GetMenuTree)       // 获取用户动态路由
		menuRouterWithoutRecord.POST("getMenuAuthority", menuApi.GetMenuAuthority) // 获取指定角色menu
		menuRouterWithoutRecord.POST("getBaseMenuById", menuApi.GetMenuById)       // 根据id获取菜单
	}
}
