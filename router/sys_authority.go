package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/middleware"
)

type AuthorityRouter struct{}

func (a *AuthorityRouter) InitAuthorityRouter(Router *gin.RouterGroup) {
	authorityRouter := Router.Group("authority").Use(middleware.OperationRecord())
	authorityRouterWithoutRecord := Router.Group("authority")
	{
		authorityRouter.POST("createAuthority", authApi.CreateAuthority)   // 创建角色
		authorityRouter.DELETE("deleteAuthority", authApi.DeleteAuthority) // 删除角色
		authorityRouter.PUT("updateAuthority", authApi.UpdateAuthority)    // 更新角色
		authorityRouter.POST("copyAuthority", authApi.CopyAuthority)       // 拷贝角色
		authorityRouter.POST("setDataAuthority", authApi.SetDataAuthority) // 设置角色资源权限
	}
	{
		authorityRouterWithoutRecord.POST("getAuthorityList", authApi.GetAuthorityList) // 获取角色列表
	}

}
