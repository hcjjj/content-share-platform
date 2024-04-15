// Package internal -----------------------------
// @file      : main.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-01-29 20:44
// -------------------------------------------
package main

func main() {
	// 初始化数据库
	//db := initDB()
	// 初始化 Redis
	//rdb := initRDB()
	// 初始化 User Handler
	//u := initUser(db, rdb)
	// 初始化 Web服务
	//server := initWebServer()
	// 注册 User 相关路由
	//u.RegisterRoutes(server)
	// 启动 Web服务
	//server.Run(":8080")

	// wire
	server := InitWebServer()
	server.Run(":8080")
}
