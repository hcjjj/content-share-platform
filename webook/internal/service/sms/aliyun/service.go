package aliyun

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	sms "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/ecodeclub/ekit"
	"github.com/goccy/go-json"
)

type Service struct {
	client *sms.Client
}

func NewService(client *sms.Client) *Service {
	return &Service{
		client: client,
	}
}

// SendSms 单次
func (s *Service) SendSms(ctx context.Context, signName, tplCode string, phone []string) error {
	phoneLen := len(phone)

	// phone1 phone2
	//    0     1
	for i := 0; i < phoneLen; i++ {
		phoneSignle := phone[i]

		// 1. 生成验证码
		code := fmt.Sprintf("%06v",
			rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

		// 完全没有做成一个独立的发短信的实现。而是一个强耦合验证码的实现。
		bcode, _ := json.Marshal(map[string]interface{}{
			"code": code,
		})

		// 2. 初始化短信结构体
		smsRequest := &sms.SendSmsRequest{
			SignName:      ekit.ToPtr[string](signName),
			TemplateCode:  ekit.ToPtr[string](tplCode),
			PhoneNumbers:  ekit.ToPtr[string](phoneSignle),
			TemplateParam: ekit.ToPtr[string](string(bcode)),
		}

		// 3. 发送短信
		smsResponse, _ := s.client.SendSms(smsRequest)
		if *smsResponse.Body.Code == "OK" {
			fmt.Println(phoneSignle, string(bcode))
			fmt.Printf("发送手机号: %s 的短信成功,验证码为【%s】\n", phoneSignle, code)
		}
		fmt.Println(errors.New(*smsResponse.Body.Message))
	}
	return nil
}
