// Package tencent -----------------------------
// @file      : service_test.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-20 12:22
// -------------------------------------------
package tencent

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func TestSender(t *testing.T) {
	secretId, ok := os.LookupEnv("SMS_SECRET_ID")
	if !ok {
		t.Fatal()
	}
	secretKey, ok := os.LookupEnv("SMS_SECRET_KEY")
	//secretId := ""
	//secretKey := ""

	c, err := sms.NewClient(common.NewCredential(secretId, secretKey),
		"ap-shanghai",
		profile.NewClientProfile())
	if err != nil {
		t.Fatal(err)
	}

	s := NewService(c, "1259389099", "hcjjj")

	testCases := []struct {
		name    string
		tplId   string
		params  []string
		numbers []string
		wantErr error
	}{
		{
			name:   "发送验证码",
			tplId:  "1877556",
			params: []string{"123456"},
			// 接收的手机号
			numbers: []string{"10086"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			er := s.Send(context.Background(), tc.tplId, tc.params, tc.numbers...)
			assert.Equal(t, tc.wantErr, er)
		})
	}
}
