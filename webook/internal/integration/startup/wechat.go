package startup

import (
	"basic-go/webook/internal/service/oauth2/wechat"
	"basic-go/webook/internal/web"
	"basic-go/webook/pkg/logger"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}

func InitWechatHandlerConfig() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{}
}
