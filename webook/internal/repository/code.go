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

type CodeRepository interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, code string) (bool, error)
}

type CachedCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CachedCodeRepository{
		cache: c,
	}
}

func (repo *CachedCodeRepository) Store(ctx context.Context, biz string,
	phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CachedCodeRepository) Verify(ctx context.Context, biz string,
	phone string, code string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, code)
}
