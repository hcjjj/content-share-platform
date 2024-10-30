// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"basic-go/webook/interaction/events"
	"basic-go/webook/interaction/grpc"
	"basic-go/webook/interaction/ioc"
	"basic-go/webook/interaction/repository"
	"basic-go/webook/interaction/repository/cache"
	"basic-go/webook/interaction/repository/dao"
	"basic-go/webook/interaction/service"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitApp() *App {
	loggerV1 := ioc.InitLogger()
	db := ioc.InitDB(loggerV1)
	interactiveDAO := dao.NewGORMInteractiveDAO(db)
	cmdable := ioc.InitRedis()
	interactiveCache := cache.NewInteractiveRedisCache(cmdable)
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	client := ioc.InitSaramaClient()
	interactiveReadEventConsumer := events.NewInteractiveReadEventConsumer(interactiveRepository, client, loggerV1)
	v := ioc.InitConsumers(interactiveReadEventConsumer)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	server := ioc.NewGrpcxServer(interactiveServiceServer, loggerV1)
	app := &App{
		consumers: v,
		server:    server,
	}
	return app
}

// wire.go:

var thirdPartySet = wire.NewSet(ioc.InitDB, ioc.InitLogger, ioc.InitSaramaClient, ioc.InitRedis)

var interactiveSvcSet = wire.NewSet(dao.NewGORMInteractiveDAO, cache.NewInteractiveRedisCache, repository.NewCachedInteractiveRepository, service.NewInteractiveService)
