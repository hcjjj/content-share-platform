package domain

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Job struct {
	Id   int64
	Name string
	// Cron 表达式
	Expression string
	Executor   string
	Cfg        string
	CancelFunc func()
}

// 通过 CORN 计算出下一次的运行时间

func (j Job) NextTime() time.Time {
	c := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := c.Parse(j.Expression)
	return s.Next(time.Now())
}
