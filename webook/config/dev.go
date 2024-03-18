//go:build !k8s

// Package config -----------------------------
// @file      : dev.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-18 11:54
// -------------------------------------------
// 没有使用 k8s 标签 就用这个的来编译
package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:13306)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:16379",
	},
}
