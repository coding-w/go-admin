package initialize

import (
	"github.com/robfig/cron/v3"
	"go-admin/global"
	"go-admin/task"
	"go.uber.org/zap"
)

func Timer() {
	go func() {
		var option []cron.Option
		// 允许秒级调度
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务 每天三点执行
		_, err := global.GA_Scheduler.AddTaskByFunc("ClearDB", "0 3 * * *", func() {
			// 定时任务方法定在task文件包中
			err := task.ClearTable(global.GA_DB)
			if err != nil {
				global.GA_LOG.Error("timer error:", zap.Error(err))
			}
		}, "定时清理数据库【日志，黑名单】内容", option...)
		if err != nil {
			global.GA_LOG.Error("add timer error:", zap.Error(err))
		}
	}()
}
