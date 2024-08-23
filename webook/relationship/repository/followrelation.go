package repository

import (
	"basic-go/webook/pkg/logger"
	"basic-go/webook/relationship/domain"
	"basic-go/webook/relationship/repository/cache"
	"basic-go/webook/relationship/repository/dao"
	"context"
)

type FollowRepository interface {
	// GetFollowee 获取某人的关注列表
	GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error)
	// FollowInfo 查看关注人的详情
	FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error)
	// AddFollowRelation 创建关注关系
	AddFollowRelation(ctx context.Context, f domain.FollowRelation) error
	// InactiveFollowRelation 取消关注
	InactiveFollowRelation(ctx context.Context, follower int64, followee int64) error
	GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error)
}

type CachedRelationRepository struct {
	dao   dao.FollowRelationDao
	cache cache.FollowCache
	l     logger.LoggerV1
}

// GetFollowStatics 获得个人的关注了多少人，以及粉丝的数量
func (d *CachedRelationRepository) GetFollowStatics(ctx context.Context, uid int64) (domain.FollowStatics, error) {
	res, err := d.cache.GetStaticsInfo(ctx, uid)
	if err == nil {
		return res, nil
	}
	// 没有就去数据库里查询
	res.Followers, err = d.dao.CntFollower(ctx, uid)
	if err != nil {
		return domain.FollowStatics{}, err
	}
	res.Followees, err = d.dao.CntFollowee(ctx, uid)
	if err != nil {
		return domain.FollowStatics{}, err
	}
	err = d.cache.SetStaticsInfo(ctx, uid, res)
	if err != nil {
		// 这里，记录一下日志
	}
	return res, nil
}

func (d *CachedRelationRepository) InactiveFollowRelation(ctx context.Context, follower int64, followee int64) error {
	err := d.dao.UpdateStatus(ctx, followee, follower, dao.FollowRelationStatusInactive)
	if err != nil {
		return err
	}
	// -1
	return d.cache.CancelFollow(ctx, follower, followee)
}

func (d *CachedRelationRepository) GetFollowee(ctx context.Context, follower, offset, limit int64) ([]domain.FollowRelation, error) {
	// 可以考虑在这里缓存关注者列表的第一页
	followerList, err := d.dao.FollowRelationList(ctx, follower, offset, limit)
	if err != nil {
		return nil, err
	}
	return d.genFollowRelationList(followerList), nil
}

func (d *CachedRelationRepository) genFollowRelationList(followerList []dao.FollowRelation) []domain.FollowRelation {
	res := make([]domain.FollowRelation, 0, len(followerList))
	for _, c := range followerList {
		res = append(res, d.toDomain(c))
	}
	return res
}

func (d *CachedRelationRepository) FollowInfo(ctx context.Context, follower int64, followee int64) (domain.FollowRelation, error) {
	// relationship:123:234 => 标签信息，分组信息
	c, err := d.dao.FollowRelationDetail(ctx, follower, followee)
	if err != nil {
		return domain.FollowRelation{}, err
	}
	return d.toDomain(c), nil
}

func (d *CachedRelationRepository) AddFollowRelation(ctx context.Context, c domain.FollowRelation) error {
	err := d.dao.CreateFollowRelation(ctx, d.toEntity(c))
	if err != nil {
		return err
	}
	// 更新缓存里面的关注了多少人，以及有多少粉丝的计数， +1
	return d.cache.Follow(ctx, c.Follower, c.Followee)
}

func (d *CachedRelationRepository) toDomain(fr dao.FollowRelation) domain.FollowRelation {
	return domain.FollowRelation{
		Followee: fr.Followee,
		Follower: fr.Follower,
	}
}

func (d *CachedRelationRepository) toEntity(c domain.FollowRelation) dao.FollowRelation {
	return dao.FollowRelation{
		Followee: c.Followee,
		Follower: c.Follower,
	}
}

func NewFollowRelationRepository(dao dao.FollowRelationDao,
	cache cache.FollowCache, l logger.LoggerV1) FollowRepository {
	return &CachedRelationRepository{
		dao:   dao,
		cache: cache,
		l:     l,
	}
}
