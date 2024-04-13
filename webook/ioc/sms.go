// Package ioc -----------------------------
// @file      : sms.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-24 19:17
// -------------------------------------------
package ioc

import (
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/localsms"
	"basic-go/webook/internal/service/sms/tencent"
	"os"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	"github.com/redis/go-redis/v9"
	tencentSMS "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	// 基于内存的实现，还是换别的
	return localsms.NewService()
	// 使用有限流的
	//return ratelimit.NewRateLimitSMSService(localsms.NewService(), limiter.NewRedisSlidingWindowLimiter(cmd, 3000, 1))
	//return initTencentSMSService()

}

func initTencentSMSService() sms.Service {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		panic("找不到腾讯 SMS 的 secret id")
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	if !ok {
		panic("找不到腾讯 SMS 的 secret key")
	}
	c, err := tencentSMS.NewClient(
		common.NewCredential(secretId, secretKey),
		"ap-nanjing",
		profile.NewClientProfile(),
	)
	if err != nil {
		panic(err)
	}
	return tencent.NewService(c, "1400842696", "妙影科技")
}
