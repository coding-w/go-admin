package router

import "github.com/gin-gonic/gin"

type DBRouter struct{}

func (db *DBRouter) InitDBRouter(Router *gin.RouterGroup) {
	initRouter := Router.Group("init")
	{
		initRouter.GET("db", dbApi.InitDB) // 初始化数据库
	}
}
