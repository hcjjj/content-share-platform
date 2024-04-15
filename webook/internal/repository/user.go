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
	"database/sql"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

type UserRepository interface {
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}
func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	// SELECT *FROM `users` WHERE `phone` = ?
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	// SELECT *FROM `users` WHERE `email` = ?
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
	// 不过好像没什么必要
	// 这边可以加个缓存 通常用户注册完就会登录 这边的过期时间可以设置得非常短
	// mail → password
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
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

	// webook/internal/repository/user_test.go
	//u = r.entityToDomain(ue)
	//_ = r.cache.Set(ctx, u)

	u = r.entityToDomain(ue)
	// 异步放入缓存
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

func (r *CachedUserRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := r.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Birthday: sql.NullInt64{
			Int64: u.Birthday.UnixMilli(),
			Valid: !u.Birthday.IsZero(),
		},
		Nickname: sql.NullString{
			String: u.Nickname,
			Valid:  u.Nickname != "",
		},
		AboutMe: sql.NullString{
			String: u.AboutMe,
			Valid:  u.AboutMe != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
		WechatOpenID: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionID: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
	}
}

func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	var birthday time.Time
	if u.Birthday.Valid {
		birthday = time.UnixMilli(u.Birthday.Int64)
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Nickname: u.Nickname.String,
		AboutMe:  u.AboutMe.String,
		Birthday: birthday,
		Ctime:    time.UnixMilli(u.Ctime),
		WechatInfo: domain.WechatInfo{
			UnionID: u.WechatUnionID.String,
			OpenID:  u.WechatOpenID.String,
		},
	}
}
