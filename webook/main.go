// Package internal -----------------------------
// @file      : main.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-01-29 20:44
// -------------------------------------------
package main

import (
	"basic-go/webook/config"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/sms/localsms"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	//db := initDB()
	// 初始化 Redis
	//rdb := initRDB()
	// 初始化 User Handler
	//u := initUser(db, rdb)
	// 初始化 Web服务
	//server := initWebServer()
	// 注册 User 相关路由
	//u.RegisterRoutes(server)
	// 启动 Web服务
	//server.Run(":8080")

	// wire
	server := InitWebServer()
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	// 限流
	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	// 1s 限流 100的请求
	// 压测的时候需要取消
	//server.Use(limiter.NewBuilder(redisClient, time.Second, 100).Build())

	// 解决跨域问题，作用于定义在这个 server 的全部路由
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods: []string{"POST", "GET"},
		// 这边需要和前端对应
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		// 加这个前端才能拿到
		ExposeHeaders: []string{"x-jwt-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境
				return true
			}
			return strings.Contains(origin, "hcjjj.webook.com")
		},
		// preflight 有效期
		MaxAge: 12 * time.Hour,
	}))

	// session management
	//store := cookie.NewStore([]byte("secret"))
	//store := memstore.NewStore([]byte("DDs0d8i62qjM8GhwhxCG3JHp6JF4Zsqc"),
	//	[]byte("vX2Vep2UjPPpr7JmMGjFcF6f0Gf8YyAc"))
	//store, err := redis.NewStore(16, "tcp", "localhost:6379",
	//	"",
	//	[]byte("DDs0d8i62qjM8GhwhxCG3JHp6JF4Zsqc"),
	//	[]byte("vX2Vep2UjPPpr7JmMGjFcF6f0Gf8YyAc"))
	//if err != nil {
	//	panic(err)
	//}
	//store := NewStore([]byte("DDs0d8i62qjM8GhwhxCG3JHp6JF4Zsqc"),
	//	[]byte("vX2Vep2UjPPpr7JmMGjFcF6f0Gf8YyAc"))
	//放 session 到每个 ctx
	//server.Use(sessions.Sessions("ssid", store))

	// 登录校验 session
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	//server.Use(middleware.CheckLogin())

	// JWT 的登录校验
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		IgnorePaths("/users/login").Build())
	return server
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	// 初始化 Uer
	// DAO repository service handler
	ud := dao.NewUserDAO(db)
	//uc := cache.NewUserCache(redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//}), time.Minute*30)
	uc := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(ud, uc)
	svc := service.NewUserService(repo)

	codeCache := cache.NewCodeCache(rdb)
	codeRepo := repository.NewCodeRepository(codeCache)
	// 方便测试 使用基于内存的 sms 实现 没有通过第三方服务
	smsSvc := localsms.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	return u
}

func initDB() *gorm.DB {
	// 打开数据库
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	// 初始化表（自动建表）
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRDB() *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
