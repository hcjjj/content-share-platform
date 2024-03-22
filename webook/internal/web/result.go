// Package web -----------------------------
// @file      : result.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-22 10:55
// -------------------------------------------
package web

type Result struct {
	// 业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
