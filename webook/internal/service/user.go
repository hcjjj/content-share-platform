// Package service -----------------------------
// @file      : user.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-02-25 11:19
// -------------------------------------------
package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("邮箱或密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type UserServiceV1 struct {
	repo repository.UserRepository
}

// NewUserService 传入的是接口 返回的是接口 为了符合 wire
func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceV1{
		repo: repo,
	}
}

func (svc *UserServiceV1) SignUp(ctx context.Context, u domain.User) error {
	// 加密 password，会影响性能
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 保存加密后的 password
	return svc.repo.Create(ctx, u)
}

func (svc *UserServiceV1) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 再比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// 打印日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserServiceV1) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	// 先查询一下这个手机号注册过没有
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		// nil 会进来
		// 用户存在也会进来
		return u, err
	}
	// 没有这个用户的话
	u = domain.User{
		Phone: phone,
	}
	// 通过新用户的手机号注册
	err = svc.repo.Create(ctx, u)
	if err != nil && err != ErrUserDuplicate {
		return u, err
	}
	// 然后再查询其 Id
	// 可能会有主从延迟的坑🕳
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *UserServiceV1) FindOrCreateByWechat(ctx context.Context,
	info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenID)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	u = domain.User{
		WechatInfo: info,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 因为这里会遇到主从延迟的问题
	return svc.repo.FindByWechat(ctx, info.OpenID)
}

func (svc *UserServiceV1) Profile(ctx context.Context, id int64) (domain.User, error) {
	// 在系统内部，基本上都是用 ID 的
	// 有些人的系统比较复杂，有一个 GUID（global unique ID）
	return svc.repo.FindById(ctx, id)
}
