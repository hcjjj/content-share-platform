package repository

import (
	"context"

	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/pkg/logger"
)

//go:generate mockgen -source=./interactive.go -package=repomocks -destination=mocks/interactive.mock.go InteractiveRepository
type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context,
		biz string, bizId int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
}

type CachedReadCntRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
	l     logger.LoggerV1
}

func (c *CachedReadCntRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	panic("implement me")
}

func (c *CachedReadCntRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	panic("implement me")
}

func (c *CachedReadCntRepository) IncrLike(ctx context.Context,
	biz string, bizId int64, uid int64) error {
	// 先插入点赞，然后更新点赞计数，更新缓存
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	// 这种做法，你需要在 repository 层面上维持住事务
	//c.dao.IncrLikeCnt()
	return c.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) DecrLike(ctx context.Context,
	biz string, bizId int64, uid int64) error {
	err := c.dao.DeleteLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return c.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) IncrReadCnt(ctx context.Context,
	biz string, bizId int64) error {
	// 要考虑缓存方案了
	// 这两个操作能不能换顺序？ —— 不能
	// MySQL
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	//go func() {
	//	c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
	//}()
	//return err
	// Redis
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) AddCollectionItem(ctx context.Context,
	biz string, bizId, cid, uid int64) error {
	panic("implement me")
}

func (c *CachedReadCntRepository) Get(ctx context.Context,
	biz string, bizId int64) (domain.Interactive, error) {
	panic("implement me")
}

func (c *CachedReadCntRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		LikeCnt:    intr.LikeCnt,
		CollectCnt: intr.CollectCnt,
		ReadCnt:    intr.ReadCnt,
	}
}

func NewCachedInteractiveRepository(dao dao.InteractiveDAO,
	cache cache.InteractiveCache, l logger.LoggerV1) InteractiveRepository {
	return &CachedReadCntRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}
