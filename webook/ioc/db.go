// Package ioc -----------------------------
// @file      : db.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 19:11
// -------------------------------------------
package ioc

import (
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/pkg/logger"
	"time"

	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	// 使用 Viper 配置 直接读
	addr := viper.GetString("db.mysql.dsn")
	// 不启动日志
	//db, err := gorm.Open(mysql.Open(addr))
	// 启动日志
	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		// 缺了一个 writer
		Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
			// 慢查询阈值，只有执行时间超过这个阈值，才会使用
			// 50ms， 100ms
			// SQL 查询必然要求命中索引，最好就是走一次磁盘 IO
			// 一次磁盘 IO 是不到 10ms
			SlowThreshold:             time.Millisecond * 10,
			IgnoreRecordNotFoundError: true,
			// 设置为 false 会直接显示具体数据
			ParameterizedQueries: true,
			LogLevel:             glogger.Info,
		}),
	})

	// 包装一下
	//type Config struct {
	//	DSN string `yaml:"dsn"`

	// 有些人的做法
	// localhost:13316
	//Addr string
	//// localhost
	//Domain string
	//// 13316
	//Port string
	//Protocol string
	//// root
	//Username string
	//// root
	//Password string
	//// webook
	//DBName string
	//}
	// 这边也可以给 cfg 一个默认值，如果配置文件没有的话
	//var cfg Config
	//err := viper.UnmarshalKey("db.mysql", &cfg)

	// remote 不支持 key 的切割 db.mysql → db 只能一级
	//err := viper.UnmarshalKey("db", &cfg)

	//if err != nil {
	//	panic(err)
	//}
	//db, err := gorm.Open(mysql.Open(cfg.DSN))

	// 使用最开始设置的配置结构体
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))

	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

// 技巧，只有单方法的接口可以这样子用

type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
