package time

import (
	cron "github.com/robfig/cron/v3"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCronExpr(t *testing.T) {
	expr := cron.New(cron.WithSeconds())
	// 这个任务的标识符
	// @every 是便利语法
	id, err := expr.AddJob("@every 1s", JobFunc(func() {
		t.Log("hello, world")
	}))
	require.NoError(t, err)
	t.Log(id)
	// 调度运行
	expr.Start()
	// 假装我们运行一段时间
	time.Sleep(time.Second * 10)
	// 停止任务
	ctx := expr.Stop()
	// 等待正在运行中的任务运行结束
	<-ctx.Done()
}

type JobFunc func()

func (jf JobFunc) Run() {
	jf()
}
