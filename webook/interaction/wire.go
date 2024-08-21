//go:build wireinject

package main

import (
	"basic-go/webook/interaction/events"
	"basic-go/webook/interaction/grpc"
	"basic-go/webook/interaction/ioc"
	repository2 "basic-go/webook/interaction/repository"
	cache2 "basic-go/webook/interaction/repository/cache"
	dao2 "basic-go/webook/interaction/repository/dao"
	service2 "basic-go/webook/interaction/service"

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
