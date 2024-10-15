package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/middleware"
)

type ApiRouter struct{}

func (ar *ApiRouter) InitApiRouter(Router *gin.RouterGroup, RouterPub *gin.RouterGroup) {
	apiRouter := Router.Group("api").Use(middleware.OperationRecord())
	apiRouterWithoutRecord := Router.Group("api")

	apiPublicRouterWithoutRecord := RouterPub.Group("api")
	{
		apiRouter.GET("getApiGroups", apiApi.GetApiGroups)          // 获取路由组
		apiRouter.GET("syncApi", apiApi.SyncApi)                    // 同步Api
		apiRouter.POST("ignoreApi", apiApi.IgnoreApi)               // 忽略Api
		apiRouter.POST("enterSyncApi", apiApi.EnterSyncApi)         // 确认同步Api
		apiRouter.POST("createApi", apiApi.CreateApi)               // 创建Api
		apiRouter.POST("deleteApi", apiApi.DeleteApi)               // 删除Api
		apiRouter.POST("getApiById", apiApi.GetApiById)             // 获取单条Api消息
		apiRouter.POST("updateApi", apiApi.UpdateApi)               // 更新api
		apiRouter.DELETE("deleteApisByIds", apiApi.DeleteApisByIds) // 删除选中api
	}
	{
		apiRouterWithoutRecord.POST("getAllApis", apiApi.GetAllApis) // 获取所有api
		apiRouterWithoutRecord.POST("getApiList", apiApi.GetApiList) // 获取Api列表
	}
	{
		apiPublicRouterWithoutRecord.GET("freshCasbin", apiApi.FreshCasbin) // 刷新casbin权限
	}
}
