// Package config -----------------------------
// @file      : k8s.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-18 11:56
// -------------------------------------------

//go:build k8s

// 使用 k8s 这个编译标签

package config

var Config = config{
	// k8s 会自动做解析
	DB: DBConfig{
		DSN: "root:root@tcp(webook-mysql:13306)/webook",
	},
	Redis: RedisConfig{
		Addr: "webook-redis:16379",
	},
}
