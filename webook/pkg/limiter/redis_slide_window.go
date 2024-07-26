package limiter

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed slide_window.lua
var luaScript string

type RedisSlidingWindowLimiter struct {
	// Redis 客户端
	cmd redis.Cmdable
	// 窗口大小
	// interval 内允许 rate 个请求
	// 1s 内允许 3000 个请求
	interval time.Duration
	// 阈值
	rate int
}

// 用 wire 所以返回接口类型
func NewRedisSlidingWindowLimiter(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSlidingWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (b *RedisSlidingWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	// 键（key）、窗口大小（interval）、速率（rate）、以及当前时间的毫秒数
	return b.cmd.Eval(ctx, luaScript, []string{key},
		b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
