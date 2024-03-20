// Package tencent -----------------------------
// @file      : service.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-20 11:07
// -------------------------------------------
package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit/slice"

	"github.com/ecodeclub/ekit"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
}

func NewService(client *sms.Client, appId string, signName string) *Service {
	return &Service{
		// string → *string
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
		client:   client,
	}
}

func (s *Service) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)

	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}

	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s", *status.Code, *status.Message)
		}
	}

	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
