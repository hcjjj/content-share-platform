// Package middleware -----------------------------
// @file      : login_jwt.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-14 11:22
// -------------------------------------------
package middleware

import (
	"basic-go/webook/internal/web"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 这边不需要登录校验
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}

		// 使用 JWT 进行校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			// 没人登录或者数据错误
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 拿出 token
		tokenStr := segs[1]
		// 用指针才能把数据塞回去
		claims := &web.UserClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("jaks3jgvkjoiGezwd4QbE9ujPZp0fL8p"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// err 为 nil token 不为 nil
		if token == nil || !token.Valid || claims.Uid == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 出现了严重的安全问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 每十秒钟刷新一次
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err := token.SignedString([]byte("jaks3jgvkjoiGezwd4QbE9ujPZp0fL8p"))
			if err != nil {
				// 记录日志
				log.Println("jwt 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Set("claims", claims)
		//ctx.Set("userId", claims.Uid)
	}
}
