package router

import (
	"github.com/gin-gonic/gin"
	"go-admin/middleware"
)

type SysDictionaryRouter struct {
}

func (s *SysDictionaryRouter) InitSysDictionaryRouter(Router *gin.RouterGroup) {
	sysDictionaryRouter := Router.Group("sysDictionary").Use(middleware.OperationRecord())
	sysDictionaryRouterWithoutRecord := Router.Group("sysDictionary")
	{
		sysDictionaryRouter.POST("createSysDictionary", dictionaryApi.CreateSysDictionary)   // 新建SysDictionary
		sysDictionaryRouter.DELETE("deleteSysDictionary", dictionaryApi.DeleteSysDictionary) // 删除SysDictionary
		sysDictionaryRouter.PUT("updateSysDictionary", dictionaryApi.UpdateSysDictionary)    // 更新SysDictionary
	}
	{
		sysDictionaryRouterWithoutRecord.GET("findSysDictionary", dictionaryApi.FindSysDictionary)       // 根据ID获取SysDictionary
		sysDictionaryRouterWithoutRecord.GET("getSysDictionaryList", dictionaryApi.GetSysDictionaryList) // 获取SysDictionary列表
	}
}
