//go:build wireinject

package startup

import (
	"basic-go/webook/interactive/repository"
	"basic-go/webook/interactive/repository/cache"
	"basic-go/webook/interactive/repository/dao"
	"basic-go/webook/interactive/service"

	"basic-go/webook/interactive/grpc"

	"github.com/google/wire"
)

// 第三方依赖
var thirdProvider = wire.NewSet(InitRedis,
	InitTestDB, InitLog)

// 服务依赖
var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitInteractiveService() service.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service.NewInteractiveService(nil, nil)
}

func InitInteractiveGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(thirdProvider, interactiveSvcProvider, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}
