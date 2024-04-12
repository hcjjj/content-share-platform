package startup

import (
	"basic-go/webook/config"

	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	//return redis.NewClient(&redis.Options{
	//	Addr: "localhost:6379",
	//})
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
