//go:build wireinject

package main

import (
	grpc2 "basic-go/webook/relationship/grpc"
	"basic-go/webook/relationship/ioc"
	"basic-go/webook/relationship/repository"
	"basic-go/webook/relationship/repository/cache"
	"basic-go/webook/relationship/repository/dao"
	"basic-go/webook/relationship/service"

	"github.com/google/wire"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMFollowRelationDAO,
	repository.NewFollowRelationRepository,
	service.NewFollowRelationService,
	grpc2.NewFollowRelationServiceServer,
	cache.NewRedisFollowCache,
)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitRedis,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		ioc.InitGrpcxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
