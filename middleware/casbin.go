package middleware

import (
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/model/common/response"
	"go-admin/service"
	"go-admin/utils"
	"strconv"
	"strings"
)

var casbinService = service.ServiceGroup.CasbinService

// CasbinHandler 拦截器
// RBAC 权限控制
func CasbinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取 JWT 声明 (claims)
		claims, _ := utils.GetClaims(c)
		// 获取请求的PATH
		path := c.Request.URL.Path
		// 获取 请求 url
		obj := strings.TrimPrefix(path, global.GA_CONFIG.System.RouterPrefix)
		// 获取请求方法
		act := c.Request.Method
		// 获取用户的角色 id
		sub := strconv.Itoa(int(claims.AuthorityId))
		// 获取 Casbin 的 SyncedCachedEnforcer 实例
		syncedCachedEnforcer := casbinService.Casbin()
		// 使用 Casbin 验证用户是否有权限访问该路径和方法
		success, _ := syncedCachedEnforcer.Enforce(sub, obj, act)
		if !success {
			response.FailWithDetailed(gin.H{}, "权限不足", c)
			c.Abort()
			return
		}
		c.Next()
	}
}
