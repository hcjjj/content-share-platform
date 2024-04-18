package ioc

import (
	"basic-go/webook/pkg/logger"

	"go.uber.org/zap"
)

// wire 依赖注入 所以需要返回接口

func InitLogger() logger.LoggerV1 {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return logger.NewZapLogger(l)
}
