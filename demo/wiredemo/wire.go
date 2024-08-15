//go:build wireinject

package wiredemo

import (
	"basic-go/demo/wiredemo/repository"
	"basic-go/demo/wiredemo/repository/dao"

	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	// 下面随便返回一个
	return &repository.UserRepository{}
}
