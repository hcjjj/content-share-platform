package dao

import "gorm.io/gorm"

type UserDAO struct {
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{}
}
