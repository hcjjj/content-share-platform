package dao

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// ErrDataNotFound 通用的数据没找到
var ErrDataNotFound = gorm.ErrRecordNotFound

//go:generate mockgen -source=./comment.go -package=daomocks -destination=mocks/comment.mock.go CommentDAO
type CommentDAO interface {
	Insert(ctx context.Context, u Comment) error
	// FindByBiz 只查找一级评论
	FindByBiz(ctx context.Context, biz string,
		bizId, minID, limit int64) ([]Comment, error)
	// FindCommentList Comment的id为0 获取一级评论，如果不为0获取对应的评论，和其评论的所有回复
	FindCommentList(ctx context.Context, u Comment) ([]Comment, error)
	FindRepliesByPid(ctx context.Context, pid int64, offset, limit int) ([]Comment, error)
	// Delete 删除本节点和其对应的子节点
	Delete(ctx context.Context, u Comment) error
	FindOneByIDs(ctx context.Context, id []int64) ([]Comment, error)
	FindRepliesByRid(ctx context.Context, rid int64, id int64, limit int64) ([]Comment, error)
}

type TreeBase struct {
	PID int64
}

// SELECT (COUNT id) FROM xx
// EXPLAIN SELECT id from xx

// Comment 把这个评论的表结构设计好
type Comment struct {
	Id int64 `gorm:"autoIncrement,primaryKey"`
	// 发表评论的人
	// 如果需要查询某个人发表的所有的评论，需要在这里创建一个索引
	Uid int64
	// 被评价的东西，建索引
	Biz     string `gorm:"index:biz_type_id"`
	BizID   int64  `gorm:"index:biz_type_id"`
	Content string
	// 根评论是哪个，也就是说，如果这个字段是 NULL，它是根评论
	RootID sql.NullInt64 `gorm:"column:root_id;index"`
	// 上一级评论的ID 这个是 NULL，也是根评论
	PID           sql.NullInt64 `gorm:"column:pid;index"`
	ParentComment *Comment      `gorm:"ForeignKey:PID;AssociationForeignKey:ID;constraint:OnDelete:CASCADE"`
	Ctime         int64
	// 事实上，大部分平台是不允许修改评论的
	Utime int64
}

func (*Comment) TableName() string {
	return "comments"
}

type GORMCommentDAO struct {
	db *gorm.DB
}

func (c *GORMCommentDAO) FindRepliesByRid(ctx context.Context,
	rid int64, id int64, limit int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("root_id = ? AND id > ?", rid, id).
		Order("id ASC").
		Limit(int(limit)).Find(&res).Error
	return res, err
}

func NewCommentDAO(db *gorm.DB) CommentDAO {
	return &GORMCommentDAO{
		db: db,
	}
}

func (c *GORMCommentDAO) FindOneByIDs(ctx context.Context, ids []int64) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).
		Where("id in ?", ids).
		First(&res).
		Error
	return res, err
}

func (c *GORMCommentDAO) FindByBiz(ctx context.Context, biz string,
	bizId, minID, limit int64) ([]Comment, error) {
	var res []Comment
	// 找的根评论
	err := c.db.WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND id < ? AND pid IS NULL", biz, bizId, minID).
		Limit(int(limit)).
		Find(&res).Error
	return res, err
}

// FindRepliesByPid 查找评论的直接评论
func (c *GORMCommentDAO) FindRepliesByPid(ctx context.Context,
	pid int64,
	offset,
	limit int) ([]Comment, error) {
	var res []Comment
	err := c.db.WithContext(ctx).Where("pid = ?", pid).
		Order("id DESC").
		Offset(offset).Limit(limit).Find(&res).Error
	return res, err
}

func (c *GORMCommentDAO) Insert(ctx context.Context, u Comment) error {
	return c.db.
		WithContext(ctx).
		Create(&u).
		Error
}

func (c *GORMCommentDAO) FindCommentList(ctx context.Context, u Comment) ([]Comment, error) {
	var res []Comment
	builder := c.db.WithContext(ctx)
	if u.Id == 0 {
		builder = builder.
			Where("biz=?", u.Biz).
			Where("biz_id=?", u.BizID).
			Where("root_id is null")
	} else {
		builder = builder.Where("root_id=? or id =?", u.Id, u.Id)
	}
	err := builder.Find(&res).Error
	return res, err

}

// 借助外键

func (c *GORMCommentDAO) Delete(ctx context.Context, u Comment) error {
	return c.db.WithContext(ctx).Delete(&Comment{
		Id: u.Id,
	}).Error
}
