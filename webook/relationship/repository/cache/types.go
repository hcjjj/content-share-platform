package cache

import (
	"basic-go/webook/relationship/domain"
	"context"
)

type FollowCache interface {
	GetStaticsInfo(ctx context.Context, uid int64) (domain.FollowStatics, error)
	SetStaticsInfo(ctx context.Context, uid int64, statics domain.FollowStatics) error
	Follow(ctx context.Context, follower, followee int64) error
	CancelFollow(ctx context.Context, follower, followee int64) error
}
