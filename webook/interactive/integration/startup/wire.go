//go:build wireinject

package startup

import (
	repository2 "basic-go/webook/interactive/repository"
	cache2 "basic-go/webook/interactive/repository/cache"
	dao2 "basic-go/webook/interactive/repository/dao"
	service2 "basic-go/webook/interactive/service"

	"github.com/google/wire"
)

var thirdProvider = wire.NewSet(InitRedis,
	InitTestDB, InitLog)
var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewCachedInteractiveRepository,
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)

func InitInteractiveService() service2.InteractiveService {
	wire.Build(thirdProvider, interactiveSvcProvider)
	return service2.NewInteractiveService(nil, nil)
}
