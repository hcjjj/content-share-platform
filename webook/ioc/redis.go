// Package ioc -----------------------------
// @file      : redis.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 19:13
// -------------------------------------------
package ioc

import (
	"github.com/spf13/viper"

	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	// 使用 Viper 配置
	addr := viper.GetString("redis.addr")
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	// 使用 配置结构体
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	return redisClient
}
