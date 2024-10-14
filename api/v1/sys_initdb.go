package v1

import (
	"github.com/gin-gonic/gin"
	"go-admin/global"
	"go-admin/model/common/response"
	"go.uber.org/zap"
)

type DBApi struct{}

func (db *DBApi) InitDB(c *gin.Context) {
	if global.GA_DB == nil {
		global.GA_LOG.Error("数据库配置错误")
		response.FailWithMessage("数据库配置错误", c)
		return
	}
	if err := initDBService.InitDB(); err != nil {
		global.GA_LOG.Error("自动创建数据库失败!", zap.Error(err))
		response.FailWithMessage("自动创建数据库失败，请查看后台日志，检查后在进行初始化", c)
		return
	}
	response.OkWithMessage("自动创建数据库成功", c)
}
