// Package middleware -----------------------------
// @file      : login.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-11 16:29
// -------------------------------------------
package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

//func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
//	l.paths = append(l.paths)
//}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//for _, path := range l.paths {
		//
		//}

		// 这些是不需要登录校验的
		if ctx.Request.URL.Path == "/users/login" ||
			ctx.Request.URL.Path == "/users/signup" {
			return
		}
		// 从 ctx 中拿出 session
		sess := sessions.Default(ctx)
		//if sess == nil {
		//	// 没有登录
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
