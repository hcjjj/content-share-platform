// Package ioc -----------------------------
// @file      : sms.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 19:17
// -------------------------------------------
package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 基于内存的实现，还是换别的
	return memory.NewService()
}
