package grpc

import (
	"context"
	"github.com/ecodeclub/ekit/queue"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"sync/atomic"
	"time"
)

type CounterLimiter struct {
	cnt       atomic.Int32
	threshold int32
}

func (c *CounterLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		// 请求进来，先占坑
		cnt := c.cnt.Add(1)
		defer func() {
			c.cnt.Add(-1)
		}()
		if cnt <= c.threshold {
			resp, err = handler(ctx, req)
			// 返回了响应
			return
		}
		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}
}

type FixedWindowLimiter struct {
	window          time.Duration
	lastWindowStart time.Time
	cnt             int
	threshold       int
	lock            sync.Mutex
}

func NewFixedWindowLimiter(window time.Duration, threshold int) *FixedWindowLimiter {
	return &FixedWindowLimiter{window: window, lastWindowStart: time.Now(), cnt: 0, threshold: threshold}

}

func (c *FixedWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (resp any, err error) {
		c.lock.Lock()
		now := time.Now()
		if now.After(c.lastWindowStart.Add(c.window)) {
			c.cnt = 0
			c.lastWindowStart = now
		}
		cnt := c.cnt + 1
		// 注意锁的范围
		c.lock.Unlock()
		if cnt <= c.threshold {
			resp, err = handler(ctx, req)
			return
		}
		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}
}

type SlidingWindowLimiter struct {
	window time.Duration
	// 请求到来的时间戳
	// 时间戳最小的在队首
	queue     queue.PriorityQueue[time.Time]
	lock      sync.Mutex
	threshold int
}

func (c *SlidingWindowLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		c.lock.Lock()
		// 我先考虑队列里面的时间戳是不是都在我的窗口范围内
		now := time.Now()

		// 快路径检测
		if c.queue.Len() < c.threshold {
			_ = c.queue.Enqueue(now)
			c.lock.Unlock()
			resp, err = handler(ctx, req)
			return
		}

		windowStart := now.Add(-c.window)
		for {
			first, _ := c.queue.Peek()
			if first.Before(windowStart) {
				// 把第一个元素删了
				_, _ = c.queue.Dequeue()
			} else {
				break
			}
		}
		if c.queue.Len() < c.threshold {
			_ = c.queue.Enqueue(now)
			c.lock.Unlock()
			resp, err = handler(ctx, req)
			return
		}
		c.lock.Unlock()
		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}
}

type TokenBucketLimiter struct {
	// 隔多久产生一个令牌
	interval  time.Duration
	buckets   chan struct{}
	closeCh   chan struct{}
	closeOnce sync.Once
}

func (c *TokenBucketLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	ticker := time.NewTicker(c.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				select {
				case c.buckets <- struct{}{}:
				default:
					// bucket 满了
				}
			case <-c.closeCh:
				return
			}
		}
	}()

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		select {
		case <-c.buckets:
			return handler(ctx, req)
		//做法1
		default:
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
			// 做法2
			//case <-ctx.Done():
			//	return nil, ctx.Err()
		}
	}
}

func (c *TokenBucketLimiter) Close() error {
	c.closeOnce.Do(func() {
		close(c.closeCh)
	})
	return nil
}

type LeakyBucketLimiter struct {
	// 隔多久产生一个令牌
	interval  time.Duration
	closeCh   chan struct{}
	closeOnce sync.Once
}

func (c *LeakyBucketLimiter) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	ticker := time.NewTicker(c.interval)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		select {
		case <-ticker.C:
			return handler(ctx, req)
		case <-c.closeCh:
			// 限流器已经关了
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
		//做法1
		default:
			return nil, status.Errorf(codes.ResourceExhausted, "限流")
			// 做法2
			//case <-ctx.Done():
			//	return nil, ctx.Err()
		}
	}
}

func (c *LeakyBucketLimiter) Close() error {
	c.closeOnce.Do(func() {
		close(c.closeCh)
	})
	return nil
}
