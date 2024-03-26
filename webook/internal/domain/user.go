// Package domain -----------------------------
// @file      : user.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-02-25 19:26
// -------------------------------------------
package domain

import "time"

// User 领域对象，是 DDD 中的 entity
// BO (Business Object)
type User struct {
	Id       int64
	Email    string
	Nickname string
	Password string
	Phone    string
	AboutMe  string
	Ctime    time.Time
	Birthday time.Time
}
