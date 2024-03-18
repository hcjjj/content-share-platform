// Package config -----------------------------
// @file      : types.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-18 11:51
// -------------------------------------------
package config

type config struct {
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	DSN string
}
type RedisConfig struct {
	Addr string
}
