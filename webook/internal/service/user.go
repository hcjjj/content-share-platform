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
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("邮箱或密码不对")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 加密 password，会影响性能
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	// 保存加密后的 password
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, uReq domain.User) (domain.User, error) {
	// 先找用户
	u, err := svc.repo.FindByEmail(ctx, uReq.Email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 再比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(uReq.Password))
	if err != nil {
		// 打印日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}
