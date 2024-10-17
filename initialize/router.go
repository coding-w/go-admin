package initialize

import (
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/middleware"
	"go-admin/router"
	"net/http"
)

// Routers 初始化路由
func Routers() *gin.Engine {
	Router := gin.New()
	// gin.Recovery()：Gin 自带的中间件，用于从 panic 中恢复，防止程序崩溃，并返回 HTTP 500 错误
	Router.Use(gin.Recovery()).Use(middleware.GinRecovery(false))
	if gin.Mode() == gin.DebugMode {
		Router.Use(gin.Logger())
	}

	systemRouter := router.RouterGroup

	// todo swagger

	PublicGroup := Router.Group(global.GA_CONFIG.System.RouterPrefix)
	PrivateGroup := Router.Group(global.GA_CONFIG.System.RouterPrefix)
	PrivateGroup.Use(middleware.AuthHandler()).Use(middleware.CasbinHandler())
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
		// 注册验证码
		systemRouter.InitCaptchaRouter(PublicGroup)
		// 初始化数据库
		systemRouter.InitDBRouter(PublicGroup)
	}
	{
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)
		systemRouter.InitAuthorityRouter(PrivateGroup)
		systemRouter.InitAuthorityBtnRouter(PrivateGroup)
		systemRouter.InitSysDictionaryRouter(PrivateGroup)
		systemRouter.InitSysMenuRouter(PrivateGroup)
		systemRouter.InitSysUserRouter(PrivateGroup, PublicGroup)
	}
	global.GA_ROUTERS = Router.Routes()
	global.GA_LOG.Info("router register success")
	initBizRouter(PrivateGroup, PublicGroup)
	return Router
}

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroup
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]

	holder(publicGroup, privateGroup)
}
