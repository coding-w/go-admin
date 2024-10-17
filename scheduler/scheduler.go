package scheduler

import (
	"github.com/robfig/cron/v3"
	"sync"
)

type Scheduler interface {
	// FindCronList 寻找所有Cron
	FindCronList() map[string]*taskManager

	// AddTaskByFuncWithSecond 添加任务函数（以秒为单位）
	AddTaskByFuncWithSecond(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error)

	// AddTaskByJobWithSeconds 添加任务接口（以秒为单位）
	AddTaskByJobWithSeconds(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error)

	// AddTaskByFunc 通过函数的方法添加任务
	AddTaskByFunc(cronName string, spec string, task func(), taskName string, option ...cron.Option) (cron.EntryID, error)

	// AddTaskByJob 通过接口的方法添加任务 要实现一个带有 Run方法的接口触发
	AddTaskByJob(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error)

	// FindCron 获取对应taskName的cron 可能会为空
	FindCron(cronName string) (*taskManager, bool)

	// StartCron 指定cron开始执行
	StartCron(cronName string)

	// StopCron 指定cron停止执行
	StopCron(cronName string)

	// FindTask 查找指定cron下的指定task
	FindTask(cronName string, taskName string) (*task, bool)

	// RemoveTask 根据id删除指定cron下的指定task
	RemoveTask(cronName string, id int)

	// RemoveTaskByName 根据taskName删除指定cron下的指定task
	RemoveTaskByName(cronName string, taskName string)

	// Clear 清理掉指定cronName
	Clear(cronName string)

	// Close 停止所有的cron
	Close()
}

// 定义定时器管理器
type schedulerManager struct {
	cronMap    map[string]*taskManager
	sync.Mutex // 并发控制
}

// FindCronList 获取所有的任务列表
func (sm *schedulerManager) FindCronList() map[string]*taskManager {
	sm.Lock()
	defer sm.Unlock()
	return sm.cronMap
}

// AddTaskByFuncWithSecond 通过函数的方法使用WithSeconds添加任务
// cronName: 需要添加任务的 Cron 名称，用于标识任务组
// spec: 任务的执行时间表达式，遵循 Cron 表达式格式
// fun: 要执行的函数，任务被触发时将调用该函数
// taskName: 任务的名称，用于在管理任务时进行查找和标识
// option: 可选参数，允许传递额外的配置选项，用于定制任务的行为
func (sm *schedulerManager) AddTaskByFuncWithSecond(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error) {
	sm.Lock()
	defer sm.Unlock()
	// 允许秒级调度
	option = append(option, cron.WithSeconds())
	cm := sm.getOrCreateCron(cronName, option...)
	id, err := cm.cron.AddFunc(spec, fun)
	cm.cron.Start()
	cm.tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

// AddTaskByJobWithSeconds 通过接口的方式使用WithSeconds添加任务
// cronName: 需要添加任务的 Cron 名称
// spec: 任务的执行时间表达式
// job: 实现了 Run() 方法的接口，用于定义任务的行为
// taskName: 任务的名称
// option: 可选参数，允许传递额外的配置选项
func (sm *schedulerManager) AddTaskByJobWithSeconds(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error) {
	sm.Lock()
	defer sm.Unlock()
	option = append(option, cron.WithSeconds())
	cm := sm.getOrCreateCron(cronName, option...)
	id, err := cm.cron.AddJob(spec, job)
	cm.cron.Start()
	cm.tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

// AddTaskByFunc 通过函数的方法添加任务
// cronName: 需要添加任务的 Cron 名称
// spec: 任务的执行时间表达式
// task: 要执行的函数
// taskName: 任务的名称
// option: 可选参数，允许传递额外的配置选项
func (sm *schedulerManager) AddTaskByFunc(cronName string, spec string, fun func(), taskName string, option ...cron.Option) (cron.EntryID, error) {
	sm.Lock()
	defer sm.Unlock()
	cm := sm.getOrCreateCron(cronName, option...)
	id, err := cm.cron.AddFunc(spec, fun)
	cm.cron.Start()
	cm.tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

// AddTaskByJob 通过接口的方法添加任务
// cronName: 需要添加任务的 Cron 名称
// spec: 任务的执行时间表达式。
// job: 实现了 Run() 方法的接口
// taskName: 任务的名称
// option: 可选参数，允许传递额外的配置选项
func (sm *schedulerManager) AddTaskByJob(cronName string, spec string, job interface{ Run() }, taskName string, option ...cron.Option) (cron.EntryID, error) {
	sm.Lock()
	defer sm.Unlock()
	cm := sm.getOrCreateCron(cronName, option...)
	id, err := cm.cron.AddJob(spec, job)
	cm.cron.Start()
	cm.tasks[id] = &task{
		EntryID:  id,
		Spec:     spec,
		TaskName: taskName,
	}
	return id, err
}

// 获取或创建指定的 Cron
func (sm *schedulerManager) getOrCreateCron(cronName string, option ...cron.Option) *taskManager {
	sm.Lock()
	defer sm.Unlock()

	if cm, exists := sm.cronMap[cronName]; exists {
		return cm
	}
	newCron := make(map[cron.EntryID]*task)
	newTaskManager := &taskManager{
		cron:  cron.New(option...),
		tasks: newCron,
	}
	sm.cronMap[cronName] = newTaskManager
	return newTaskManager
}

// FindCron 查找指定名称的 Cron
// cronName: 要查找的 Cron 名称。
func (sm *schedulerManager) FindCron(cronName string) (*taskManager, bool) {
	sm.Lock()
	defer sm.Unlock()
	cm, ok := sm.cronMap[cronName]
	return cm, ok
}

// StartCron 启动指定的 Cron
// cronName: 要启动的 Cron 名称
func (sm *schedulerManager) StartCron(cronName string) {
	sm.Lock()
	defer sm.Unlock()
	if v, ok := sm.cronMap[cronName]; ok {
		v.cron.Start()
	}
}

// StopCron 停止指定的 Cron
// cronName: 要停止的 Cron 名称
func (sm *schedulerManager) StopCron(cronName string) {
	sm.Lock()
	defer sm.Unlock()
	if v, ok := sm.cronMap[cronName]; ok {
		v.cron.Stop()
	}
}

// FindTask 查找指定 Cron 下的任务
// cronName: 要查找的 Cron 名称
// taskName: 要查找的任务名称
func (sm *schedulerManager) FindTask(cronName string, taskName string) (*task, bool) {
	sm.Lock()
	defer sm.Unlock()
	manager, ok := sm.cronMap[cronName]
	if !ok {
		return nil, false
	}
	for _, t2 := range manager.tasks {
		if t2.TaskName == taskName {
			return t2, true
		}
	}
	return nil, false
}

// RemoveTask 根据 ID 删除指定 Cron 下的任务
// cronName: 要删除任务的 Cron 名称
// id: 任务的 EntryID
func (sm *schedulerManager) RemoveTask(cronName string, id int) {
	sm.Lock()
	defer sm.Unlock()
	if v, ok := sm.cronMap[cronName]; ok {
		v.cron.Remove(cron.EntryID(id))
		delete(v.tasks, cron.EntryID(id))
	}
}

// RemoveTaskByName 根据任务名称删除指定 Cron 下的任务
// cronName: 要删除任务的 Cron 名称
// taskName: 要删除的任务名称
func (sm *schedulerManager) RemoveTaskByName(cronName string, taskName string) {
	fTask, ok := sm.FindTask(cronName, taskName)
	if !ok {
		return
	}
	sm.RemoveTask(cronName, int(fTask.EntryID))
}

// Clear 清理掉指定的 Cron
// cronName: 要清理的 Cron 名称
func (sm *schedulerManager) Clear(cronName string) {
	sm.Lock()
	defer sm.Unlock()
	v, ok := sm.cronMap[cronName]
	if ok {
		v.cron.Stop()
		delete(sm.cronMap, cronName)
	}
}

// Close 停止所有Cron
func (sm *schedulerManager) Close() {
	sm.Lock()
	defer sm.Unlock()
	for _, v := range sm.cronMap {
		v.cron.Stop()
	}
	sm.cronMap = make(map[string]*taskManager)
}

func NewSchedulerManager() Scheduler {
	return &schedulerManager{cronMap: make(map[string]*taskManager)}
}
