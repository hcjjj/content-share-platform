// Package ioc -----------------------------
// @file      : db.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 19:11
// -------------------------------------------
package ioc

import (
	dao2 "basic-go/webook/interactive/repository/dao"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/pkg/logger"
	"time"

	promsdk "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

func InitDB(l logger.LoggerV1) *gorm.DB {

	// 使用 Viper 配置 直接读
	addr := viper.GetString("db.mysql.dsn")

	// 不启动日志
	db, err := gorm.Open(mysql.Open(addr))

	// 启动日志
	//db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
	//	// 缺了一个 writer
	//	Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
	//		// 慢查询阈值，只有执行时间超过这个阈值，才会使用
	//		// 50ms， 100ms
	//		// SQL 查询必然要求命中索引，最好就是走一次磁盘 IO
	//		// 一次磁盘 IO 是不到 10ms
	//		SlowThreshold:             time.Millisecond * 10,
	//		IgnoreRecordNotFoundError: true,
	//		// 设置为 false 会直接显示具体数据
	//		ParameterizedQueries: true,
	//		LogLevel:             glogger.Info,
	//	}),
	//})

	// 使用 Prometheus 监控插件
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}
	// 使用 callback 统计执行时间
	// 监控查询的执行时间
	pcb := newCallbacks()
	//pcb.registerAll(db)
	// 插件的用法 自己包装的
	db.Use(pcb)

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
	// 数据库表初始化
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	err = dao2.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

type Callbacks struct {
	vector *promsdk.SummaryVec
}

func (pcb *Callbacks) Name() string {
	return "prometheus-query"
}

func (pcb *Callbacks) Initialize(db *gorm.DB) error {
	pcb.registerAll(db)
	return nil
}

func newCallbacks() *Callbacks {
	vector := promsdk.NewSummaryVec(promsdk.SummaryOpts{
		// 在这边，你要考虑设置各种 Namespace
		Namespace: "hcjjj",
		Subsystem: "webook",
		Name:      "gorm_query_time",
		Help:      "统计 GORM 的执行时间",
		ConstLabels: map[string]string{
			"db": "webook",
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	},
		// 如果是 JOIN 查询，table 就是 JOIN 在一起的
		// 或者 table 就是主表，A JOIN B，记录的是 A
		[]string{"type", "table"})

	pcb := &Callbacks{
		vector: vector,
	}
	// 上报 Prometheus
	promsdk.MustRegister(vector)
	return pcb
}

func (pcb *Callbacks) registerAll(db *gorm.DB) {
	// 作用于 INSERT 语句
	err := db.Callback().Create().Before("*").
		Register("prometheus_create_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").
		Register("prometheus_create_after", pcb.after("create"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Update().Before("*").
		Register("prometheus_update_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").
		Register("prometheus_update_after", pcb.after("update"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Delete().Before("*").
		Register("prometheus_delete_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").
		Register("prometheus_delete_after", pcb.after("delete"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Raw().Before("*").
		Register("prometheus_raw_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").
		Register("prometheus_raw_after", pcb.after("raw"))
	if err != nil {
		panic(err)
	}

	err = db.Callback().Row().Before("*").
		Register("prometheus_row_before", pcb.before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").
		Register("prometheus_row_after", pcb.after("row"))
	if err != nil {
		panic(err)
	}
}

func (c *Callbacks) before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}

func (c *Callbacks) after(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			// 啥都干不了
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknown"
		}
		c.vector.WithLabelValues(typ, table).
			Observe(float64(time.Since(startTime).Milliseconds()))
	}
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
