package ioc

import (
	"basic-go/webook/comment/repository/dao"
	"basic-go/webook/pkg/logger"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	c := Config{
		DSN: "root:root@tcp(localhost:3306)/mysql",
	}
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", c, err))
	}
	db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
		//使用 DEBUG 来打印
		Logger: glogger.Default.LogMode(glogger.Info),
	})
	if err != nil {
		panic(err)
	}

	// 接入 prometheus
	err = db.Use(prometheus.New(prometheus.Config{
		DBName: "webook",
		// 每 15 秒采集一些数据
		RefreshInterval: 15,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"Threads_running"},
			},
		}, // user defined metrics
	}))
	if err != nil {
		panic(err)
	}
	err = db.Use(tracing.NewPlugin(tracing.WithoutMetrics()))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
