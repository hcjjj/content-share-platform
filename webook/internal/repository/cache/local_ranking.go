package cache

import (
	"basic-go/webook/internal/domain"
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"time"
)

type RankingLocalCache struct {
	// 用泛型封装
	// 可以考虑直接使用 uber 的，或者 SDK 自带的
	topN       *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func NewRankingLocalCache() *RankingLocalCache {
	return &RankingLocalCache{
		topN: atomicx.NewValue[[]domain.Article](),
		ddl:  atomicx.NewValueOf(time.Now()),
		// 永不过期，或者非常长，或者对齐到 redis 的过期时间，都行
		expiration: time.Minute * 10,
	}
}

//func (r *RankingLocalCache) Preload(ctx context.Context) {
//
//}

func (r *RankingLocalCache) Set(ctx context.Context, arts []domain.Article) error {
	// 也可以按照 id => Article 缓存
	r.topN.Store(arts)
	ddl := time.Now().Add(r.expiration)
	r.ddl.Store(ddl)
	return nil
}

func (r *RankingLocalCache) Get(ctx context.Context) ([]domain.Article, error) {
	ddl := r.ddl.Load()
	arts := r.topN.Load()
	if len(arts) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("本地缓存未命中")
	}
	return arts, nil
}
func (r *RankingLocalCache) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := r.topN.Load()
	return arts, nil
}

type item struct {
	arts []domain.Article
	ddl  time.Time
}
