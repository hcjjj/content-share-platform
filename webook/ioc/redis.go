package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	addr := viper.GetString("redis.addr")
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	return redisClient
}

func InitRLockClient(cmd redis.Cmdable) *rlock.Client {
	return rlock.NewClient(cmd)
}

//
//func NewRateLimiter() redis.Limiter {
//
//}
