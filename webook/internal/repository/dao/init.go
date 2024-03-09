// Package dao -----------------------------
// @file      : init.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-02-25 20:29
// -------------------------------------------
package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) error {
	// 自动建表
	return db.AutoMigrate(&User{})
}
