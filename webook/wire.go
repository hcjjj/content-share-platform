// Package main -----------------------------
// @file      : wiredemo.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 18:46
// -------------------------------------------

//go:build wireinject

package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		// 业务相关
		dao.NewUserDAO,
		cache.NewUserCache, cache.NewCodeCache,
		repository.NewUserRepository, repository.NewCodeRepository,
		ioc.InitSMSService,
		service.NewUserService, service.NewCodeService,
		web.NewUserHandler,
		// 微信登录
		//ioc.InitWechatService,
		//ioc.NewWechatHandlerConfig,
		//web.NewOAuth2WechatHandler,
		// 中间件
		ioc.InitMiddlewares,
		// web（服务 + 路由）
		ioc.InitWebServer,
	)
	return new(gin.Engine)
}
