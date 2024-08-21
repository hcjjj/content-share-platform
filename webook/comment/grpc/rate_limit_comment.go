package grpc

import (
	commentv1 "basic-go/webook/api/proto/gen/comment/v1"
	"context"
	"errors"
)

type RateLimitComment struct {
	CommentServiceServer
}

func (c *RateLimitComment) GetMoreReplies(ctx context.Context, req *commentv1.GetMoreRepliesRequest) (*commentv1.GetMoreRepliesResponse, error) {
	if ctx.Value("limited") == "true" || ctx.Value("downgrade") == "true" {
		return &commentv1.GetMoreRepliesResponse{}, errors.New("资源不足，功能关闭")
	}
	return c.CommentServiceServer.GetMoreReplies(ctx, req)
}

func (c *RateLimitComment) GetCommentList(ctx context.Context, request *commentv1.CommentListRequest) (*commentv1.CommentListResponse, error) {
	// 一般是通过热榜功能，提前计算放到了 Redis 里面，问一下 Redis 就知道是不是热门资源了
	isHotBiz := c.isHotBiz(request.Biz, request.Bizid)
	if isHotBiz {
		// 限流阈值 400/s
	} else {
		// 限流阈值 100/s
	}
	return c.CommentServiceServer.GetCommentList(ctx, request)
}

func (c *RateLimitComment) GetCommentListV1(ctx context.Context, request *commentv1.CommentListRequest) (*commentv1.CommentListResponse, error) {
	// 一般是通过热榜功能，提前计算放到了 Redis 里面，问一下 Redis 就知道是不是热门资源了
	isHotBiz := c.isHotBiz(request.Biz, request.Bizid)
	if !isHotBiz && ctx.Value("downgrade") == "true" {
		// 非热门资源，触发降级
		return &commentv1.CommentListResponse{}, errors.New("非热门资源被降级")
	}
	return c.CommentServiceServer.GetCommentList(ctx, request)
}

func (c *RateLimitComment) CreateComment(ctx context.Context, request *commentv1.CreateCommentRequest) (*commentv1.CreateCommentResponse, error) {
	if ctx.Value("limited") == "true" || ctx.Value("downgrade") == "true" {
		// 转 Kafka
		return &commentv1.CreateCommentResponse{}, nil
	}
	err := c.svc.CreateComment(ctx, convertToDomain(request.GetComment()))
	return &commentv1.CreateCommentResponse{}, err
}

func (c *RateLimitComment) isHotBiz(biz string, bizid int64) bool {
	return true
}
