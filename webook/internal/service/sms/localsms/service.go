// Package localsms -----------------------------
// @file      : limiter.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-22 10:26
// -------------------------------------------
package localsms

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
func (s Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	// 模拟发短信过程
	fmt.Printf("短信验证码测试：%s\n", args)
	return nil
}
