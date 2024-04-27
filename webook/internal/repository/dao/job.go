package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (g *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id =?", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).
		Where("id = ?", id).Updates(map[string]any{
		"status": jobStatusPaused,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) Release(ctx context.Context, id int64) error {
	// 这里有一个问题。你要不要检测 status 或者 version?
	// WHERE version = ?
	// 要 记得修改 防止自己挂机后任务在别人那 然后把别人的任务给释放了
	return g.db.WithContext(ctx).Model(&Job{}).Where("id =?", id).
		Updates(map[string]any{
			"status": jobStatusWaiting,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	// 高并发情况下，大部分都是陪太子读书
	// 100 个 goroutine
	// 要转几次？ 所有 goroutine 执行的循环次数加在一起是
	// 1+2+3+4 +5 + ... + 99 + 100
	// 特定一个 goroutine，最差情况下，要循环一百次

	// 超时控制
	db := g.db.WithContext(ctx)
	for {
		now := time.Now()
		var j Job
		// 分布式任务调度系统 方案 高并发环境下 一般没有高并发的情况的
		// 1. 一次拉一批，我一次性取出 100 条来，然后，我随机从某一条开始，向后开始抢占
		// 2. 我搞个随机偏移量，0-100 生成一个随机偏移量。兜底：第一轮没查到，偏移量回归到 0
		// 3. 我搞一个 id 取余分配，status = ? AND next_time <=? AND id%10 = ? 兜底：不加余数条件，取next_time 最老的
		err := db.WithContext(ctx).Where("status = ? AND next_time <=?", jobStatusWaiting, now).
			First(&j).Error
		// 你找到了，可以被抢占的
		// 找到之后你要干嘛？你要抢占
		if err != nil {
			// 没有任务。从这里返回
			return Job{}, err
		}
		// 两个 goroutine 都拿到 id =1 的数据
		// 能不能用 utime?
		// 乐观锁，CAS 操作，compare AND Swap
		// 有一个很常见的面试刷亮点：就是用乐观锁取代 FOR UPDATE
		// 性能优化：曾用 FOR UPDATE =>性能差，还会有死锁 => 优化成了乐观锁
		res := db.Where("id=? AND version = ?",
			j.Id, j.Version).Model(&Job{}).
			Updates(map[string]any{
				"status": jobStatusRunning,
				"utime":  now,
				// 使用乐观锁方案
				"version": j.Version + 1,
			})
		if res.Error != nil {
			return Job{}, err
		}
		if res.RowsAffected == 0 {
			// 抢占失败，要继续下一轮
			continue
		}
		return j, nil
	}
}

//

type Job struct {
	Id       int64 `gorm:"primaryKey,autoIncrement"`
	Cfg      string
	Executor string
	Name     string `gorm:"unique"`

	// 第一个问题：哪些任务可以抢？哪些任务已经被人占着？哪些任务永远不会被运行
	// 用状态来标记
	Status int

	// 另外一个问题，定时任务，我怎么知道，已经到时间了呢？
	// NextTime 下一次被调度的时间
	// next_time <= now 这样一个查询条件
	// and status = 0
	// 要建立索引
	// 更加好的应该是 next_time 和 status 的联合索引
	NextTime int64 `gorm:"index"`
	// cron 表达式
	Cron string

	// 解决并发问题
	Version int

	// 创建时间，毫秒数
	Ctime int64
	// 更新时间，毫秒数
	Utime int64
}

const (
	jobStatusWaiting = iota
	// 已经被抢占
	jobStatusRunning
	// 还可以有别的取值

	// 暂停调度
	jobStatusPaused
)
