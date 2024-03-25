// Package cache -----------------------------
// @file      : code.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-20 19:38
// -------------------------------------------
package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var (
	ErrSetCodeTooMany         = errors.New("发送验证码太频繁")
	ErrVerifyCodeTooManyTimes = errors.New("验证次数超过限制")
	ErrUnknownForCode         = errors.New("位置错误 与 code 有关")
)

// 编译器会在编译的时候，把 set_code 的代码放进来这个 luaSetCode 变量里
//
//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	key(biz, phone string) string
	Verify(ctx context.Context, biz, phone, expectedCode string) (bool, error)
}

type RedisCodeCache struct {
	client redis.Cmdable
}

// NewCodeCacheGoBestPractice Go 的最佳实践是返回具体类型
func NewCodeCacheGoBestPractice(client redis.Cmdable) *RedisCodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

// NewCodeCache wire 需要返回接口
func NewCodeCache(client redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		client: client,
	}
}

func (c *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		// 毫无问题
		return nil
	case -1:
		// 发送太频繁
		return ErrSetCodeTooMany
	default:
		// 系统错误
		return errors.New("系统错误")
	}
}

func (c *RedisCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (c *RedisCodeCache) Verify(ctx context.Context, biz, phone, expectedCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, expectedCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrVerifyCodeTooManyTimes
	case -2:
		return false, nil
		//default:
		//	return false, ErrUnknownForCode
	}
	return false, ErrUnknownForCode
}
