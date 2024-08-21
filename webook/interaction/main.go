package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	// 读取配置文件信息
	initViper()
	// 初始化服务
	app := InitApp()
	// 启动消费者
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	// 启动服务
	err := app.server.Serve()
	log.Println(err)
}

//func mainV1() {
//	initViper()
//	server := grpc2.NewServer()
//	// 这里暂时随便搞一下
//	// 搞成依赖注入
//	// 这种写法的缺陷是，如果我有很多个 grpc API 服务端的实现
//	intrSvc := InitGRPCServer()
//	intrv1.RegisterInteractiveServiceServer(server, intrSvc)
//	// 监听 8090 端口，你可以随便写
//	l, err := net.Listen("tcp", ":8090")
//	if err != nil {
//		panic(err)
//	}
//	// 这边会阻塞，类似与 gin.Run
//	err = server.Serve(l)
//	log.Println(err)
//}

func initViper() {
	//--config=config/dev.yaml
	cfile := pflag.String("config",
		"./config/dev.yaml", "指定配置文件路径")
	pflag.Parse()
	viper.SetConfigFile(*cfile)

	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉你文件变了，不能告诉你，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉你变更前的数据，和变更后的数据
		// 更好的设计是，它会直接告诉你差异。
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
