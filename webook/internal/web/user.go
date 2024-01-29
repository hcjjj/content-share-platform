// Package web -----------------------------
// @file      : user.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-01-29 20:45
// -------------------------------------------
package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 定义用户有关的路由
type UserHandler struct {
}

func (u *UserHandler) RegisterRoutesV1(ug *gin.RouterGroup) {
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	//server.POST("/users/signup", u.SignUp)
	//server.POST("/users/login", u.Login)
	//server.POST("/users/edit", u.Edit)
	//server.GET("/users/profile", u.Profile)
	ug := server.Group("/users")
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	// Bind 方法会根据 Content-type 解析数据到 req
	// 解析出错回直接写回一个 400 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%v\n", req)
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}
func (u *UserHandler) Profile(ctx *gin.Context) {

}
