package grpc

import (
	"basic-go/webook/pkg/limiter"
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LimiterUserServer struct {
	limiter limiter.Limiter
	UserServiceServer
}

func (s *LimiterUserServer) GetByID(ctx context.Context,
	req *GetByIDRequest) (*GetByIDResponse, error) {
	key := fmt.Sprintf("limiter:user:get_by_id:%d", req.Id)
	limited, err := s.limiter.Limit(ctx, key)
	if err != nil {
		// 有保守的做法，也有激进的做法
		// 这个是保守的做法
		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}

	if limited {
		return nil, status.Errorf(codes.ResourceExhausted, "限流")
	}
	return s.UserServiceServer.GetByID(ctx, req)
}
