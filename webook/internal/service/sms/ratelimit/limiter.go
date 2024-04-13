package ratelimit

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/pkg/limiter"
	"context"
	"errors"
)

var errLimited = errors.New("触发限流")

var _ sms.Service = &RateLimitSMSService{}

type RateLimitSMSService struct {
	// 被装饰的
	// 不使用组合的方式，需要自己去实现 sms.Service 的所有方法
	svc     sms.Service
	limiter limiter.Limiter
	key     string
}

type RateLimitSMSServiceV1 struct {
	// 使用组合的方式
	// 自动实现 sms.Service 的所有方法
	// 只需要装饰我需要装饰的方法即可
	sms.Service
	limiter limiter.Limiter
	key     string
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		// 系统错误
		// 可以限流：保守策略，下游很脆的时候
		// 可以不限：下游很强的时候，容错策略
		return err
	}
	if limited {
		return errLimited
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}

func NewRateLimitSMSService(svc sms.Service,
	l limiter.Limiter) *RateLimitSMSService {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: l,
		key:     "sms-limiter",
	}
}
