package aliyun

import (
	"context"
	"testing"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	sms "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/ecodeclub/ekit"
	"github.com/stretchr/testify/assert"
)

//func TestSender(t *testing.T) {
//
//	keyId := ""
//	keySecret := ""
//
//	config := &openapi.Config{
//		AccessKeyId:     ekit.ToPtr[string](keyId),
//		AccessKeySecret: ekit.ToPtr[string](keySecret),
//	}
//	client, err := sms.NewClient(config)
//	if err != nil {
//		t.Fatal(err)
//	}
//	service := NewService(client)
//
//	testCases := []struct {
//		signName string
//		tplCode  string
//		phone    string
//		wantErr  error
//	}{
//		{
//			signName: "webook",
//			tplCode:  "SMS_462745194",
//			phone:    "",
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.signName, func(t *testing.T) {
//			er := service.SendSms(context.Background(), tc.signName, tc.tplCode, tc.phone)
//			assert.Equal(t, tc.wantErr, er)
//		})
//	}
//}

func TestService_SendSms(t *testing.T) {

	keyId := ""
	keySecret := ""

	config := &openapi.Config{
		AccessKeyId:     ekit.ToPtr[string](keyId),
		AccessKeySecret: ekit.ToPtr[string](keySecret),
	}
	client, err := sms.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(client)

	tests := []struct {
		signName string
		tplCode  string
		phone    []string
		wantErr  error
	}{
		{
			signName: "",
			tplCode:  "",
			phone:    []string{"", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.signName, func(t *testing.T) {
			er := service.SendSms(context.Background(), tt.signName, tt.tplCode, tt.phone)
			assert.Equal(t, tt.wantErr, er)
		})
	}
}
