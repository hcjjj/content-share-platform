package main

import (
	"basic-go/webook/pkg/grpcx"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	initViperV2Watch()
	app := Init()
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}

func initViperV2Watch() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	pflag.Parse()
	// 直接指定文件路径
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

type App struct {
	server *grpcx.Server
}
