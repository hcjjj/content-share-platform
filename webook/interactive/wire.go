//go:build wireinject

package main

import (
	"basic-go/webook/interactive/events"
	"basic-go/webook/interactive/grpc"
	"basic-go/webook/interactive/ioc"
	repository2 "basic-go/webook/interactive/repository"
	cache2 "basic-go/webook/interactive/repository/cache"
	dao2 "basic-go/webook/interactive/repository/dao"
	service2 "basic-go/webook/interactive/service"

	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitSaramaClient,
	ioc.InitRedis,
)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)

func InitApp() *App {
	wire.Build(thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitConsumers,
		ioc.NewGrpcxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
