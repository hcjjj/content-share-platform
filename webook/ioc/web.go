// Package ioc -----------------------------
// @file      : web.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 19:24
// -------------------------------------------
package ioc

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/ginx/middlewares/ratelimit"
	"basic-go/webook/pkg/limiter"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/redis/go-redis/v9"

	"github.com/gin-gonic/gin"
)

// 带微信登录的
//func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, oauth2WechatHdl *web.OAuth2WechatHandler) *gin.Engine {
//	server := gin.Default()
//	server.Use(mdls...)
//	userHdl.RegisterRoutes(server)
//  oauth2WechatHdl.RegisterRoutes(server)
//	return server
//}

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHlf(),
		ratelimitHlf(redisClient),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/users/login").Build(),
	}
}

func corsHlf() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:2999"},
		AllowMethods: []string{"POST", "GET"},
		// 这边需要和前端对应
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		// 加这个前端才能拿到
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				// 开发环境
				return true
			}
			return strings.Contains(origin, "hcjjj.webook.com")
		},
		// preflight 有效期
		MaxAge: 11 * time.Hour,
	})
}

func ratelimitHlf(redisClient redis.Cmdable) gin.HandlerFunc {
	return ratelimit.NewBuilder(limiter.NewRedisSlidingWindowLimiter(redisClient, time.Second, 100)).Build()
}
