// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"basic-go/webook/relationship/grpc"
	"basic-go/webook/relationship/repository"
	"basic-go/webook/relationship/repository/cache"
	"basic-go/webook/relationship/repository/dao"
	"basic-go/webook/relationship/service"
)

// Injectors from wire.go:

func InitServer() *grpc.FollowServiceServer {
	gormDB := InitTestDB()
	followRelationDao := dao.NewGORMFollowRelationDAO(gormDB)
	cmdable := InitRedis()
	followCache := cache.NewRedisFollowCache(cmdable)
	loggerV1 := InitLog()
	followRepository := repository.NewFollowRelationRepository(followRelationDao, followCache, loggerV1)
	followRelationService := service.NewFollowRelationService(followRepository)
	followServiceServer := grpc.NewFollowRelationServiceServer(followRelationService)
	return followServiceServer
}
