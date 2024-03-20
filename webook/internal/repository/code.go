// Package repository -----------------------------
// @file      : code.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-20 19:51
// -------------------------------------------
package repository

import (
	"basic-go/webook/internal/repository/cache"
	"context"
)

var (
	ErrSendCodeTooMany        = cache.ErrSetCodeTooMany
	ErrVerifyCodeTooManyTimes = cache.ErrVerifyCodeTooManyTimes
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz string,
	phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz string,
	phone string, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}
