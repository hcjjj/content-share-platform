// Package middleware -----------------------------
// @file      : login.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-11 16:29
// -------------------------------------------
package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	// 用 go 的方式编码解码
	gob.Register(time.Now())

	return func(ctx *gin.Context) {
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		// 从 ctx 中拿出 session
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 刷新登录状态（session的有效期）
		// 固定间隔时间刷新，比如说每分钟内第一次访问我都刷新
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60 * 10,
		})
		now := time.Now()
		if updateTime == nil {
			// 刚登录还没刷新过，第一次登录的第一个请求
			sess.Set("update_time", now)
			//sess.Save()
			if err := sess.Save(); err != nil {
				panic(err)
			}
			return
		}
		// updateTime 是有的
		updateTimeVal, ok := updateTime.(time.Time)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		// 超过 1 min了， 刷新
		if now.Sub(updateTimeVal) > time.Minute*1 {
			sess.Set("update_time", now)
			if err := sess.Save(); err != nil {
				panic(err)
			}
		}
	}
}

// CheckLogin 这种没上面的方式好
func CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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
