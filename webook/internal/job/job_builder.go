package job

import (
	"strconv"
	"time"

	"basic-go/webook/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

type CronJobBuilder struct {
	l logger.LoggerV1
	p *prometheus.SummaryVec
	//tracer trace.Tracer
}

func NewCronJobBuilder(l logger.LoggerV1) *CronJobBuilder {
	p := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "hcjjj",
		Subsystem: "webook",
		Help:      "统计定时任务的执行情况",
		Name:      "cron_job",
	}, []string{"name", "success"})
	prometheus.MustRegister(p)
	return &CronJobBuilder{
		l: l,
		p: p,
		//tracer: otel.GetTracerProvider().Tracer("webook/internal/job"),
	}
}

func (b *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobFuncAdapter(func() error {
		//_, span := b.tracer.Start(context.Background(), name)
		//defer span.End()
		start := time.Now()
		b.l.Info("任务开始",
			logger.String("job", name))
		var success bool
		defer func() {
			b.l.Info("任务结束",
				logger.String("job", name))
			duration := time.Since(start).Milliseconds()
			b.p.WithLabelValues(name,
				strconv.FormatBool(success)).Observe(float64(duration))
		}()
		err := job.Run()
		success = err == nil
		if err != nil {
			//span.RecordError(err)
			b.l.Error("运行任务失败", logger.Error(err),
				logger.String("job", name))
		}
		return nil
	})
}

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
	_ = c()
}
