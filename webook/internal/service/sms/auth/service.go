package auth

import (
	"basic-go/webook/internal/service/sms"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type SMSService struct {
	svc sms.Service
	key []byte
}

func (s *SMSService) Send(ctx context.Context, tplToken string, args []string, numbers ...string) error {
	var claims SMSClaims
	// 如果能解析成功，说明是最有的业务方
	_, err := jwt.ParseWithClaims(tplToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(ctx, claims.Tpl, args, numbers...)
}

type SMSClaims struct {
	jwt.RegisteredClaims
	Tpl string
	// 额外加字段
}
