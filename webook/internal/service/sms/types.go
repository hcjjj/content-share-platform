// Package sms -----------------------------
// @file      : types.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-20 10:21
// -------------------------------------------
package sms

import (
	"context"
)

type Service interface {
	Send(ctx context.Context, tplId string, args []string, numbers ...string) error
}
