//go:build wireinject

package main

import (
	"basic-go/webook/interactive/events"
	"basic-go/webook/interactive/grpc"
	"basic-go/webook/interactive/ioc"
	"basic-go/webook/interactive/repository"
	"basic-go/webook/interactive/repository/cache"
	"basic-go/webook/interactive/repository/dao"
	"basic-go/webook/interactive/service"

	"github.com/google/wire"
)

var thirdPartySet = wire.NewSet(ioc.InitDB,
	ioc.InitLogger,
	ioc.InitKafka,
	// 暂时不理会 consumer 怎么启动
	ioc.InitRedis)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitAPP() *App {
	wire.Build(interactiveSvcProvider,
		thirdPartySet,
		events.NewInteractiveReadEventConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.NewConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
