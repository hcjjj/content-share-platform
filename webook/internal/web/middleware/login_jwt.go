// Package middleware -----------------------------
// @file      : login_jwt.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-14 11:22
// -------------------------------------------
package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"

	ijwt "basic-go/webook/internal/web/jwt"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}

}

// IgnorePaths 中间方法，构建部分
func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

// Build 完成构建，返回最终需要的
func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 Go 的方式编码解码
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		tokenStr := l.ExtractToken(ctx)
		claims := &ijwt.UserClaims{}
		// ParseWithClaims 里面，一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("jaks3jgvkjoiGezwd4QbE9ujPZp0fL8p"), nil
		})
		if err != nil {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//claims.ExpiresAt.Time.Before(time.Now()) {
		//	// 过期了
		//}
		// err 为 nil，token 不为 nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			// 你是要监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 检查 Token 是否在黑名单
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// 要么 redis 有问题，要么已经退出登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 你以为的退出登录，没有用的
		//token.Valid = false
		//// tokenStr 是一个新的字符串
		//tokenStr, err = token.SignedString([]byte("jaks3jgvkjoiGezwd4QbE9ujPZp0fL8p"))
		//if err != nil {
		//	// 记录日志
		//	log.Println("jwt 续约失败", err)
		//}
		//ctx.Header("x-jwt-token", tokenStr)

		// 短的 token 过期了，搞个新的
		//now := time.Now()
		// 每十秒钟刷新一次
		//if claims.ExpiresAt.Sub(now) < time.Second*50 {
		//	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
		//	tokenStr, err = token.SignedString([]byte("jaks3jgvkjoiGezwd4QbE9ujPZp0fL8p"))
		//	if err != nil {
		//		// 记录日志
		//		log.Println("jwt 续约失败", err)
		//	}
		//	ctx.Header("x-jwt-token", tokenStr)
		//}
		ctx.Set("claims", claims)
		//ctx.Set("userId", claims.Uid)
	}
}
