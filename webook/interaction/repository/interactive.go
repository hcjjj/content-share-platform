package repository

import (
	"basic-go/webook/interaction/domain"
	"basic-go/webook/interaction/repository/cache"
	"basic-go/webook/interaction/repository/dao"
	"basic-go/webook/pkg/logger"
	"context"
	"github.com/ecodeclub/ekit/slice"
)

//go:generate mockgen -source=./interactive.go -package=repomocks -destination=mocks/interactive.mock.go InteractiveRepository
type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context,
		biz string, bizId int64) error
	// BatchIncrReadCnt 这里调用者要保证 bizs 和 bizIds 长度一样
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error
	IncrLike(ctx context.Context, biz string, bizId, uid int64) error
	DecrLike(ctx context.Context, biz string, bizId, uid int64) error
	AddCollectionItem(ctx context.Context, biz string, bizId, cid int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
	Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error)
}

type CachedReadCntRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDAO
	l     logger.LoggerV1
}

func (c *CachedReadCntRepository) GetByIds(ctx context.Context, biz string, ids []int64) ([]domain.Interactive, error) {
	vals, err := c.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	return slice.Map[dao.Interactive, domain.Interactive](vals,
		func(idx int, src dao.Interactive) domain.Interactive {
			return c.toDomain(src)
		}), nil
}

func (c *CachedReadCntRepository) Liked(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetLikeInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrDataNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedReadCntRepository) Collected(ctx context.Context, biz string, id int64, uid int64) (bool, error) {
	_, err := c.dao.GetCollectionInfo(ctx, biz, id, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrDataNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (c *CachedReadCntRepository) IncrLike(ctx context.Context,
	biz string, bizId int64, uid int64) error {
	err := c.dao.InsertLikeInfo(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
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
	err := c.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	// 这边会有部分失败引起的不一致的问题，但是你其实不需要解决，
	// 因为阅读数不准确完全没有问题
	return c.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) BatchIncrReadCnt(ctx context.Context,
	bizs []string, bizIds []int64) error {
	return c.dao.BatchIncrReadCnt(ctx, bizs, bizIds)
}

func (c *CachedReadCntRepository) AddCollectionItem(ctx context.Context,
	biz string, bizId, cid, uid int64) error {
	err := c.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
		Biz:   biz,
		Cid:   cid,
		BizId: bizId,
		Uid:   uid,
	})
	if err != nil {
		return err
	}
	return c.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (c *CachedReadCntRepository) Get(ctx context.Context,
	biz string, bizId int64) (domain.Interactive, error) {
	intr, err := c.cache.Get(ctx, biz, bizId)
	if err == nil {
		// 缓存只缓存了具体的数字，但是没有缓存自身有没有点赞的信息
		// 因为一个人反复刷，重复刷一篇文章是小概率的事情
		// 也就是说，你缓存了某个用户是否点赞的数据，命中率会很低
		return intr, nil
	}
	ie, err := c.dao.Get(ctx, biz, bizId)
	if err == nil {
		res := c.toDomain(ie)
		if er := c.cache.Set(ctx, biz, bizId, res); er != nil {
			c.l.Error("回写缓存失败",
				logger.Int64("bizId", bizId),
				logger.String("biz", biz),
				logger.Error(er))
		}
		return res, nil
	}
	return domain.Interactive{}, err
}

func (c *CachedReadCntRepository) toDomain(intr dao.Interactive) domain.Interactive {
	return domain.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
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
