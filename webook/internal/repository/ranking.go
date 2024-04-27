package repository

import (
	"context"

	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	// 使用具体实现，可读性更好，对测试不友好，因为没有面向接口编程
	redis *cache.RankingRedisCache
	local *cache.RankingLocalCache
}

func (c *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	// 先本地再Redis
	data, err := c.local.Get(ctx)
	if err == nil {
		return data, nil
	}
	data, err = c.redis.Get(ctx)
	if err == nil {
		c.local.Set(ctx, data)
	} else {
		return c.local.ForceGet(ctx)
	}
	return data, err
}

func NewCachedRankingRepository(
	redis *cache.RankingRedisCache,
	local *cache.RankingLocalCache,
) RankingRepository {
	return &CachedRankingRepository{local: local, redis: redis}
}

func (c *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	// 先本地再Redis
	_ = c.local.Set(ctx, arts)
	return c.redis.Set(ctx, arts)
}
