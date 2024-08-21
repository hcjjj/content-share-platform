package service

import (
	"basic-go/webook/comment/domain"
	"basic-go/webook/comment/repository"
	"context"
)

type CommentService interface {
	// GetCommentList Comment的id为0 获取一级评论
	// 按照 ID 倒序排序
	GetCommentList(ctx context.Context, biz string, bizId, minID, limit int64) ([]domain.Comment, error)
	// DeleteComment 删除评论，删除本评论何其子评论
	DeleteComment(ctx context.Context, id int64) error
	// CreateComment 创建评论
	CreateComment(ctx context.Context, comment domain.Comment) error
	GetMoreReplies(ctx context.Context, rid int64, maxID int64, limit int64) ([]domain.Comment, error)
}

type commentService struct {
	repo repository.CommentRepository
}

func (c *commentService) GetMoreReplies(ctx context.Context,
	rid int64,
	maxID int64, limit int64) ([]domain.Comment, error) {
	return c.repo.GetMoreReplies(ctx, rid, maxID, limit)
}

func NewCommentSvc(repo repository.CommentRepository) CommentService {
	return &commentService{
		repo: repo,
	}
}

func (c *commentService) GetCommentList(ctx context.Context, biz string,
	bizId, minID, limit int64) ([]domain.Comment, error) {
	list, err := c.repo.FindByBiz(ctx, biz, bizId, minID, limit)
	if err != nil {
		return nil, err
	}
	return list, err
}

func (c *commentService) DeleteComment(ctx context.Context, id int64) error {
	return c.repo.DeleteComment(ctx, domain.Comment{
		Id: id,
	})
}

func (c *commentService) CreateComment(ctx context.Context, comment domain.Comment) error {
	return c.repo.CreateComment(ctx, comment)
}
