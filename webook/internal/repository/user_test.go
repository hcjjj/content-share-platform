// Package repository -----------------------------
// @file      : user_test.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-04-12 10:39
// -------------------------------------------
package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	cachemocks "basic-go/webook/internal/repository/cache/mocks"
	"basic-go/webook/internal/repository/dao"
	daomocks "basic-go/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	now := time.Now()
	// 去掉毫秒意外的部分
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存未命中，查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "password",
						Phone: sql.NullString{
							String: "18105050505",
							Valid:  true,
						},
						Ctime: now.UnixMilli(),
						Utime: now.UnixMilli(),
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "password",
					Phone:    "18105050505",
					Ctime:    now,
				}).Return(nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "password",
				Phone:    "18105050505",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中，查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存未命中，查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						Password: "password",
						Phone:    "18105050505",
						Ctime:    now,
					}, nil)
				d := daomocks.NewMockUserDAO(ctrl)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "password",
				Phone:    "18105050505",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				// 缓存未命中，查了缓存，但是没结果
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("mock db 错误"))
				return d, c
			},
			ctx:      context.Background(),
			id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("mock db 错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewUserRepository(ud, uc)
			u, err := repo.FindById(tc.ctx, tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
			// 测试异步的时候可以这样子
			time.Sleep(time.Second)
		})
	}
}
