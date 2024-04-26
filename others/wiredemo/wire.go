//go:build wireinject

package wiredemo

import (
	"basic-go/others/wiredemo/repository"
	"basic-go/others/wiredemo/repository/dao"

	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	// 下面随便返回一个
	return &repository.UserRepository{}
}
