// Package internal -----------------------------
// @file      : main.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-01-29 20:44
// -------------------------------------------
package main

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	db := initDB()
	// 初始化 Web服务
	server := initWebServer()
	// 初始化 User Handler
	u := initUser(db)
	// 注册 User 相关路由
	u.RegisterRoutes(server)
	// 启动 Web服务
	server.Run(":8080")
}

func initWebServer() *gin.Engine {
	server := gin.Default()
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
			return strings.Contains(origin, "yourcompany.com")
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
	// 放 session 到每个 ctx
	//server.Use(sessions.Sessions("ssid", store))

	// 登录校验 session
	//server.Use(middleware.NewLoginMiddlewareBuilder().
	//	IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	//server.Use(middleware.CheckLogin())

	// JWT 的登录校验
	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
		IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())
	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	// 初始化 Uer
	// DAO repository service handler
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}

func initDB() *gorm.DB {
	// 打开数据库
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"))
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
