package startup

import (
	intrv1 "basic-go/webook/api/proto/gen/interaction/v1"
	"basic-go/webook/interaction/service"
	"basic-go/webook/internal/client"
	"context"

	"google.golang.org/grpc"
)

func InitIntrClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	return client.NewLocalInteractiveServiceAdapter(svc)
}

type DoNothingInteractiveServiceClient struct {
}

func (d *DoNothingInteractiveServiceClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	return &intrv1.IncrReadCntResponse{}, nil
}

func (d *DoNothingInteractiveServiceClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return &intrv1.LikeResponse{}, nil
}

func (d *DoNothingInteractiveServiceClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	return &intrv1.CancelLikeResponse{}, nil
}

func (d *DoNothingInteractiveServiceClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return &intrv1.CollectResponse{}, nil
}

func (d *DoNothingInteractiveServiceClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	return &intrv1.GetResponse{
		Intr: &intrv1.Interactive{},
	}, nil
}

func (d *DoNothingInteractiveServiceClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	return &intrv1.GetByIdsResponse{
		Intrs: map[int64]*intrv1.Interactive{},
	}, nil
}
