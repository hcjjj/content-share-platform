// Package internal -----------------------------
// @file      : main.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-01-29 20:44
// -------------------------------------------
package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.uber.org/zap"

	"github.com/spf13/pflag"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	// 需要用到这里面的一个初始化方法，所以需要匿名引入
	// 使用 viper 接入 etcd
	_ "github.com/spf13/viper/remote"
)

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

	// 配置模块
	//initViper()
	initViperWithArg()

	//initPrometheus()

	// 要先把数据存在 etcd
	//initViperRemote()
	//keys := viper.AllKeys()
	//println(keys)
	//setting := viper.AllSettings()
	//fmt.Println(setting)

	// 日志模块
	initLogger()

	// wire
	app := InitWebServer()

	// 消费者
	// Consumer 的设计，类似于 Web，或者 GRPC 之类的，是一个顶级入口
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	// 定时任务
	app.cron.Start()
	// HTTP 服务
	server := app.web
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":8080")
	// 服务关闭
	// 一分钟内要关完，要退出
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ctx = app.cron.Stop()
	// 这边可以考虑超时强制退出，防止有些任务，执行特别长的时间
	tm := time.NewTimer(time.Minute * 10)
	select {
	case <-tm.C:
	case <-ctx.Done():
	}
}

func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("这是 replace 之前")
	// 如果不 replace，直接用 zap.L()，啥都打不出来。
	zap.ReplaceGlobals(logger)
	//zap.L().Info("zap start!")

	//type Demo struct {
	//	Name string `json:"name"`
	//}
	//zap.L().Info("这是实验参数",
	//	zap.Error(errors.New("这是一个 error")),
	//	zap.Int64("id", 123),
	//	zap.Any("一个结构体", Demo{Name: "hello"}))
}

func initViper() {
	// 如果配置文件里面没有 设置默认值
	viper.SetDefault("mysql.dsn",
		"root:root@tcp(localhost:13306)/mysql")
	// 配置文件的名字，但是不包含文件扩展名
	// 不包含 .go, .yaml 之类的后缀
	viper.SetConfigName("dev")
	// 告诉 viper 配置用的是 yaml 格式
	// 现实中，有很多格式，JSON，XML，YAML，TOML，ini
	viper.SetConfigType("yaml")
	// 当前工作目录下的 config 子目录
	viper.AddConfigPath("./webook/config")
	//viper.AddConfigPath("/tmp/config")
	//viper.AddConfigPath("/etc/webook")
	// 读取配置到 viper 里面，或者可以理解为加载到内存里面
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// 可以有多个 Viper 实例
	//otherViper := viper.New()
	//otherViper.SetConfigName("myjson")
	//otherViper.AddConfigPath("./config")
	//otherViper.SetConfigType("json")
}

func initViperWatch() {
	cfile := pflag.String("config",
		"config/config.yaml", "配置文件路径")
	// 这一步之后，cfile 里面才有值
	pflag.Parse()
	//viper.Set("db.dsn", "localhost:3306")
	// 所有的默认值放好s
	viper.SetConfigType("yaml")
	viper.SetConfigFile(*cfile)
	viper.WatchConfig()
	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	val := viper.Get("test.key")
	log.Println(val)
}

func initViperWithArg() {
	// 利用 viper 读取启动参数
	// 不同环境加载不同配置环境
	// --config=config/dev.yaml
	// program argument
	// go run . --config=config/dev.yaml
	cfile := pflag.String("config",
		"./webook/config/dev.yaml", "指定配置文件路径")
	// 顺序不能乱 要先从命令行中解析出来，不然都是默认值
	pflag.Parse()
	viper.SetConfigFile(*cfile)

	// 实时监听配置变更
	viper.WatchConfig()
	// 只能告诉文件变了，不能告诉，文件的哪些内容变了
	viper.OnConfigChange(func(in fsnotify.Event) {
		// 比较好的设计，它会在 in 里面告诉变更前的数据，和变更后的数据
		// 更好的设计是，它会直接告诉差异。
		fmt.Println(in.Name, in.Op)
		fmt.Println(viper.GetString("db.dsn"))
	})

	//viper.SetDefault("db.mysql.dsn",
	//	"root:root@tcp(localhost:3306)/mysql")
	//viper.SetConfigFile("config/dev.yaml")
	//viper.KeyDelimiter("-")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

// 用于开发联调环境
// 都手写了，你干嘛不直接在 IoC 里面写死呢？
func initViperReader() {
	viper.SetConfigType("yaml")
	cfg := `
db.mysql:
  dsn: "root:root@tcp(localhost:13306)/webook"

redis:
  addr: "localhost:16379"
`
	err := viper.ReadConfig(bytes.NewReader([]byte(cfg)))
	if err != nil {
		panic(err)
	}
}

func initViperRemote() {
	err := viper.AddRemoteProvider("etcd3",
		// 通过 webook 和其他使用 etcd 的区别出来
		"http://127.0.0.1:12379", "/webook")
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	//err = viper.WatchRemoteConfig()
	//if err != nil {
	//	panic(err)
	//}
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	fmt.Println(in.Name, in.Op)
	//})
	err = viper.ReadRemoteConfig()
	if err != nil {
		panic(err)
	}
}
