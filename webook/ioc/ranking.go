package ioc

import (
	"time"

	"basic-go/webook/internal/job"
	"basic-go/webook/internal/service"
	"basic-go/webook/pkg/logger"

	rlock "github.com/gotomicro/redis-lock"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(svc service.RankingService,
	rlockClient *rlock.Client,
	l logger.LoggerV1) *job.RankingJob {
	return job.NewRankingJob(svc, rlockClient, l, time.Second*30)
}

func InitJobs(l logger.LoggerV1, rankingJob *job.RankingJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	cbd := job.NewCronJobBuilder(l)
	// 这里每三分钟一次
	_, err := res.AddJob("0 */3 * * * ?", cbd.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return res
}
