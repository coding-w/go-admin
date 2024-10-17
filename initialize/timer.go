package initialize

//
//import (
//	"fmt"
//	"github.com/flipped-aurora/gin-vue-admin/server/task"
//
//	"github.com/flipped-aurora/gin-vue-admin/server/global"
//)
//
//func Timer() {
//	go func() {
//		var option []cron.Option
//		option = append(option, cron.WithSeconds())
//		timer := timer.NewTimerTask()
//		// 清理DB定时任务
//		_, err := timer.AddTaskByFunc("ClearDB", "@daily", func() {
//			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
//			if err != nil {
//				fmt.Println("timer error:", err)
//			}
//		}, "定时清理数据库【日志，黑名单】内容", option...)
//		if err != nil {
//			fmt.Println("add timer error:", err)
//		}
//	}()
//}
