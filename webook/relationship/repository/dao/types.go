package dao

import "context"

// 存储用户的关注数据
type FollowRelation struct {
	ID int64 `gorm:"column:id;autoIncrement;primaryKey;"`

	// 要在这两个列上，创建一个联合唯一索引
	// 如果认为查询一个人关注了多少人，是主要查询场景
	// <follower, followee>
	// 如果认为查询一个人有哪些粉丝，是主要查询场景
	// <followee, follower>
	// 查我关注了哪些人？ WHERE follower = 123(我的 uid)
	Follower int64 `gorm:"uniqueIndex:follower_followee"`
	Followee int64 `gorm:"uniqueIndex:follower_followee"`

	// 软删除策略
	Status uint8

	// 如果关注有类型，有优先级，有一些备注数据的
	// Type string
	// Priority string
	// Gid 分组ID

	Ctime int64
	Utime int64
}

const (
	FollowRelationStatusUnknown uint8 = iota
	FollowRelationStatusActive
	FollowRelationStatusInactive
)

type FollowRelationDao interface {
	// FollowRelationList 获取某人的关注列表
	FollowRelationList(ctx context.Context, follower, offset, limit int64) ([]FollowRelation, error)
	FollowRelationDetail(ctx context.Context, follower int64, followee int64) (FollowRelation, error)
	// CreateFollowRelation 创建联系人
	CreateFollowRelation(ctx context.Context, c FollowRelation) error
	// UpdateStatus 更新状态
	UpdateStatus(ctx context.Context, followee int64, follower int64, status uint8) error
	// CntFollower 统计计算关注自己的人有多少
	CntFollower(ctx context.Context, uid int64) (int64, error)
	// CntFollowee 统计自己关注了多少人
	CntFollowee(ctx context.Context, uid int64) (int64, error)
}

// UserRelation 另外一种设计方案，但是不要这么做
type UserRelation struct {
	ID     int64 `gorm:"primaryKey,autoIncrement,column:id"`
	Uid1   int64 `gorm:"column:uid1;type:int(11);not null;uniqueIndex:user_contact_index"`
	Uid2   int64 `gorm:"column:uid2;type:int(11);not null;uniqueIndex:user_contact_index"`
	Block  bool  // 拉黑
	Mute   bool  // 屏蔽
	Follow bool  // 关注
}

type UserRelationV1 struct {
	ID   int64 `gorm:"primaryKey,autoIncrement,column:id"`
	Uid1 int64 `gorm:"column:uid1;type:int(11);not null;uniqueIndex:user_contact_index"`
	Uid2 int64 `gorm:"column:uid2;type:int(11);not null;uniqueIndex:user_contact_index"`
	Type string
}

type FollowStatics struct {
	ID  int64 `gorm:"primaryKey,autoIncrement,column:id"`
	Uid int64 `gorm:"unique"`
	// 有多少粉丝
	Followers int64
	// 关注了多少人
	Followees int64

	Utime int64
	Ctime int64
}
