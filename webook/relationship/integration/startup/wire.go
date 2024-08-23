//go:build wireinject

package startup

import (
	"basic-go/webook/relationship/grpc"
	"basic-go/webook/relationship/repository"
	"basic-go/webook/relationship/repository/cache"
	"basic-go/webook/relationship/repository/dao"
	"basic-go/webook/relationship/service"

	"github.com/google/wire"
)

func InitServer() *grpc.FollowServiceServer {
	wire.Build(
		InitRedis,
		InitTestDB,
		InitLog,
		dao.NewGORMFollowRelationDAO,
		cache.NewRedisFollowCache,
		repository.NewFollowRelationRepository,
		service.NewFollowRelationService,
		grpc.NewFollowRelationServiceServer,
	)
	return new(grpc.FollowServiceServer)
}
