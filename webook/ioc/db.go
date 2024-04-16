// Package ioc -----------------------------
// @file      : db.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 19:11
// -------------------------------------------
package ioc

import (
	"basic-go/webook/internal/repository/dao"

	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	// 使用 Viper 配置 直接读
	addr := viper.GetString("db.mysql.dsn")
	db, err := gorm.Open(mysql.Open(addr))

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
