package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/middleware"
)

type SysUserRouter struct {
}

func (sur *SysUserRouter) InitSysUserRouter(PrivateRouter *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	userRouterPrivate := PrivateRouter.Group("/user").Use(middleware.OperationRecord())
	userRouterPrivateWithoutRecord := PrivateRouter.Group("/user")
	userRouterPublic := PublicRouter.Group("/user")

	{
		userRouterPublic.POST("/login", userApi.Login)
	}
	{
		userRouterPrivate.POST("/admin_register", userApi.Register)
		userRouterPrivate.POST("/changePassword", userApi.ChangePassword)
		userRouterPrivate.POST("/setUserAuthority", userApi.SetUserAuthority)
		userRouterPrivate.DELETE("/deleteUser", userApi.DeleteUser)
		userRouterPrivate.PUT("/setUserInfo", userApi.SetUserInfo)
		userRouterPrivate.PUT("/setSelfInfo", userApi.SetSelfInfo)
		userRouterPrivate.POST("/setUserAuthorities", userApi.SetUserAuthorities)
		userRouterPrivate.POST("/resetPassword", userApi.ResetPassword)
	}
	{
		userRouterPrivateWithoutRecord.POST("/getUserList", userApi.GetUserList)
		userRouterPrivateWithoutRecord.GET("/getUserInfo", userApi.GetUserInfo)
	}

}
