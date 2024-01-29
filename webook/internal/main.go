// Package internal -----------------------------
// @file      : main.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-01-29 20:44
// -------------------------------------------
package main

import (
	"basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	u := &web.UserHandler{}
	u.RegisterRoutes(server)
	//u.RegisterRoutesV1(server.Group())

	server.Run(":8080")
}
