package startup

import (
	"basic-go/webook/pkg/logger"
)

func InitLog() logger.LoggerV1 {
	return logger.NewNopLogger()
}
