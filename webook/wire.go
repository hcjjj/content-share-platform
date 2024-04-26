// Package main -----------------------------
// @file      : wiredemo.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 18:46
// -------------------------------------------

//go:build wireinject

package main

import (
	articleEve "basic-go/webook/internal/events/article"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/article"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	articleDao "basic-go/webook/internal/repository/dao/article"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/ioc"

	"github.com/google/wire"
)

func InitWebServer() *App {
	wire.Build(
		// 最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		// 业务相关
		// 用户
		dao.NewUserDAO,
		cache.NewUserCache, cache.NewCodeCache,
		repository.NewUserRepository, repository.NewCodeRepository,
		ioc.InitSMSService,
		service.NewUserService, service.NewCodeService,
		web.NewUserHandler,
		ijwt.NewRedisJWTHandler,
		// kafka
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,
		// 文章
		cache.NewRedisArticleCache,
		cache.NewRedisInteractiveCache,
		articleDao.NewGORMArticleDAO,
		dao.NewGORMInteractiveDAO,
		article.NewArticleRepository,
		repository.NewCachedInteractiveRepository,
		articleEve.NewKafkaProducer,
		articleEve.NewInteractiveReadEventBatchConsumer,
		service.NewArticleService,
		web.NewArticleHandler,
		// 日志模块
		ioc.InitLogger,
		// 微信登录
		//ioc.InitWechatService,
		//ioc.NewWechatHandlerConfig,
		//web.NewOAuth2WechatHandler,
		// 中间件
		ioc.InitMiddlewares,
		// web（服务 + 路由）
		ioc.InitWebServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
