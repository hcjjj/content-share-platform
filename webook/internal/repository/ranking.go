package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"context"
)

type RankingRepository interface {
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type CachedRankingRepository struct {
	cache cache.RankingCache

	// 下面是给 v1 用的
	redisCache *cache.RankingRedisCache
	localCache *cache.RankingLocalCache
}

func NewCachedRankingRepositoryV1(redisCache *cache.RankingRedisCache, localCache *cache.RankingLocalCache) *CachedRankingRepository {
	return &CachedRankingRepository{redisCache: redisCache, localCache: localCache}

}

func (repo *CachedRankingRepository) GetTopNV1(ctx context.Context) ([]domain.Article, error) {
	res, err := repo.localCache.Get(ctx)
	if err == nil {
		return res, nil
	}
	res, err = repo.redisCache.Get(ctx)
	if err != nil {
		return repo.localCache.ForceGet(ctx)
	}
	_ = repo.localCache.Set(ctx, res)
	return res, nil
}

func (repo *CachedRankingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return repo.cache.Get(ctx)
}

func NewCachedRankingRepository(cache cache.RankingCache) RankingRepository {
	return &CachedRankingRepository{cache: cache}
}

func (repo *CachedRankingRepository) ReplaceTopNV1(ctx context.Context, arts []domain.Article) error {
	_ = repo.localCache.Set(ctx, arts)
	return repo.redisCache.Set(ctx, arts)
}

func (repo *CachedRankingRepository) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return repo.cache.Set(ctx, arts)
}
