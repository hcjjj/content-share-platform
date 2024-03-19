// Package repository -----------------------------
// @file      : user.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-02-25 19:35
// -------------------------------------------
package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: c,
	}
}
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// SELECT *FROM `users` WHERE `email` = ?
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *UserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	// 不过好像没什么必要
	// 这边可以加个缓存 通常用户注册完就会登录 这边的过期时间可以设置得非常短
	// mail → password
}

func (r *UserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 从 cache 中找
	u, err := r.cache.Get(ctx, id)
	// 缓存里有数据
	if err == nil {
		return u, nil
	}

	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
	// 放入缓存
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 打日志做监控
		}
	}()
	return u, err

	// 缓存里面没有数据
	//if err == cache.ErrKeyNotExist {
	// 去数据库里面找
	//}
	// 缓存出错
	// 如果是 Redis 挂掉了 直接冲击到数据库
	// 选加载数据库需要做好兜底（限流）

}
