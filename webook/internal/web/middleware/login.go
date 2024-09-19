package middleware

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
}

// CheckLogin 登录校验实现（刷新登录态）
func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// Gin 的 session 机制使用了 GOB 来将对象转化为字节切片[]byte，所以需要提前注册一下 Gob
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}

		sess := sessions.Default(ctx)
		userId := sess.Get("userId")
		if userId == nil {
			// 中断，不要往后执行，也就是不要执行后面的业务逻辑
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		//ctx.Next()// 执行业务
		// 在执行业务之后搞点什么
		//duration := time.Now().Sub(now)

		// 固定时间间隔刷新策略
		// 假如说，策略是每分钟刷一次，怎么知道，已经过了一分钟
		const updateTimeKey = "update_time"
		// 试着拿出上一次刷新时间
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)

		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			// 这是第一次进来
			sess.Set(updateTimeKey, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				// 打日志
				fmt.Println(err)
			}
		}
	}
}
