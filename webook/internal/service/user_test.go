// Package service -----------------------------
// @file      : user_test.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-03-29 10:27
// -------------------------------------------
package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	repomocks "basic-go/webook/internal/repository/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"
)

func TestUserServiceV1_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		// 测试准备
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository
		// 输入数据
		//ctx      context.Context
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name:     "登录成功",
			email:    "123@qq.com",
			password: "hello#world",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email: "123@qq.com",
					// 这个密码是加密后的
					Password: "$2a$10$guS0GDPUGIQBGMwgk4P.OOI2B.WcBk3TGqHBsRPSSr3K0oYPgQeeK",
					Phone:    "18101234231",
					Ctime:    now,
				}, nil)
				return repo
			},

			wantUser: domain.User{
				Email: "123@qq.com",
				// 这个密码是加密后的
				Password: "$2a$10$guS0GDPUGIQBGMwgk4P.OOI2B.WcBk3TGqHBsRPSSr3K0oYPgQeeK",
				Phone:    "18101234231",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("mock db 错误"))
				return repo
			},
			email:    "123@qq.com",
			password: "hello#world123",

			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
		{
			name:     "密码错误",
			email:    "123@qq.com",
			password: "xxxhello#world",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{
					Email: "123@qq.com",
					// 这个密码是加密后的
					Password: "$2a$10$guS0GDPUGIQBGMwgk4P.OOI2B.WcBk3TGqHBsRPSSr3K0oYPgQeeK",
					Phone:    "18101234231",
					Ctime:    now,
				}, nil)
				return repo
			},

			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			// 具体的测试代码
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
