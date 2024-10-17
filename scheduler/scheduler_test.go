package scheduler

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockJob struct {
}

var job = mockJob{}

func (job mockJob) Run() {
	mockFun()
}

func mockFun() {
	fmt.Println("mockFun")
	fmt.Println("finish")
}

func TestNewSchedulerManager(t *testing.T) {
	sm := NewSchedulerManager()
	_sm := sm.(*schedulerManager)

	{
		_, err := sm.AddTaskByFunc("func", "@every 1s", mockFun, "测试mockFunc")
		assert.Nil(t, err)
		_, ok := _sm.cronMap["func"]
		if !ok {
			t.Error("no find func")
		}
	}

	{
		id, err := sm.AddTaskByJob("job", "0 3 * * *", job, "测试job mockfunc")
		assert.Nil(t, err)
		cron, ok := _sm.cronMap["job"]
		if !ok {
			t.Error("no find job")
		} else {
			t2 := cron.tasks[id]
			fmt.Println(t2)
		}
	}

	{
		_, ok := sm.FindCron("func")
		if !ok {
			t.Error("no find func")
		}
		_, ok = sm.FindCron("job")
		if !ok {
			t.Error("no find job")
		}
		_, ok = sm.FindCron("none")
		if ok {
			t.Error("find none")
		}
	}
	{
		sm.Clear("func")
		_, ok := sm.FindCron("func")
		if ok {
			t.Error("find func")
		}
	}
	{
		a := sm.FindCronList()
		b, c := sm.FindCron("job")
		fmt.Println(a, b, c)
	}
}
