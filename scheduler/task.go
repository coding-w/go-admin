package scheduler

import "github.com/robfig/cron/v3"

// 定义一个任务结构
type task struct {
	EntryID  cron.EntryID // 任务的 ID
	Spec     string       // Cron 表达式
	TaskName string       // 任务名称
}

// 定义任务管理器
type taskManager struct {
	cron  *cron.Cron             // Cron 实例
	tasks map[cron.EntryID]*task // 存储任务的映射
}
