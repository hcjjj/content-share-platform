// Package service -----------------------------
// @file      : code.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-20 17:33
// -------------------------------------------
package service

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/service/sms"
	"context"
	"fmt"
	"math/rand"
)

var (
	ErrCodeVerifyTooManyTimes = repository.ErrVerifyCodeTooManyTimes
	ErrCodeSendTooMany        = repository.ErrSendCodeTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type CodeServiceV1 struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &CodeServiceV1{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

// Send 发送验证码
func (svc *CodeServiceV1) Send(ctx context.Context,
	// 区别验证码的业务场景
	biz string,
	phone string) error {
	// 生成一个验证码
	code := svc.generateCode()
	// 存入 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 发送验证码
	// [api 不可用 这边直接把验证码发回作为流程测试]
	err = svc.smsSvc.Send(ctx, "123", []string{code}, phone)

	//if err != nil {
	// redis 有验证码 但是没有发送出去
	// 能不能删掉这个验证码
	// 如果是超时的err呢 都不知道有没有发出去
	//}
	return err
}

func (svc *CodeServiceV1) SendTest(ctx context.Context,
	// 区别验证码的业务场景
	biz string,
	phone string) (error, string) {
	// 生成一个验证码
	code := svc.generateCode()
	// 存入 Redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err, "-1"
	}
	// 发送验证码
	return nil, code
}

func (svc *CodeServiceV1) generateCode() string {
	// 0 - 9999
	num := rand.Intn(10000)
	// 补全为 4 位
	return fmt.Sprintf("%04d", num)
}

func (svc *CodeServiceV1) Verify(ctx context.Context, biz string,
	phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}
