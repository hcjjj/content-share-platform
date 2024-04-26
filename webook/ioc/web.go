// Package ioc -----------------------------
// @file      : web.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 19:24
// -------------------------------------------
package ioc

import (
	"basic-go/webook/internal/web"
	ijwt "basic-go/webook/internal/web/jwt"
	"basic-go/webook/internal/web/middleware"
	"basic-go/webook/pkg/ginx"
	"basic-go/webook/pkg/ginx/middlewares/metric"
	"basic-go/webook/pkg/ginx/middlewares/ratelimit"
	"basic-go/webook/pkg/limiter"
	logger2 "basic-go/webook/pkg/logger"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"

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

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	articleHdl.RegisterRoutes(server)
	// 观察测试
	(&web.ObservabilityHandler{}).RegisterRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable,
	l logger2.LoggerV1,
	jwtHdl ijwt.Handler) []gin.HandlerFunc {

	//bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
	//	// 这边传入打什么级别的，与打印格式
	//	l.Debug("HTTP请求", logger2.Field{Key: "al", Value: al})
	//}).AllowReqBody(true).AllowRespBody()
	//// 监听配置文件的变动，事实控制
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	ok := viper.GetBool("web.logreq")
	//	bd.AllowReqBody(ok)
	//})

	// Prometheus 监控
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "hcjjj",
		Subsystem: "webook",
		Name:      "http_biz_code",
		Help:      "HTTP 的业务错误码",
	})

	return []gin.HandlerFunc{
		// 跨域
		corsHlf(),
		// HTTP日志记录
		//bd.Build(),
		// IP 限流
		ratelimitHlf(redisClient),
		// Prometheus 监控
		(&metric.MiddlewareBuilder{
			Namespace:  "hcjjj",
			Subsystem:  "webook",
			Name:       "gin_http",
			Help:       "统计 GIN 的 HTTP 接口",
			InstanceID: "my-instance-1",
		}).Build(),

		// 不校验
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/oauth2/wechat/authurl").
			IgnorePaths("/oauth2/wechat/callback").
			IgnorePaths("/test/metric").
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
