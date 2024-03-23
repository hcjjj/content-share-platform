// Package dao -----------------------------
// @file      : user.go
// @author    : hcjjj
// @contact   : hcjjj@foxmail.com
// @time      : 2024-02-25 19:39
// -------------------------------------------
package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"

	"gorm.io/gorm"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}
func (dao *UserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("`id` = ?", id).First(&u).Error
	return u, err
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	// 存毫秒数
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	// 判断是 mysql 的错误 和底层强耦合的代码
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		// MySQL唯一键冲突错误码
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			// 邮箱冲突（因为只有邮箱是唯一索引）
			return ErrUserDuplicateEmail
		}
	}
	// 为什么不先查询来判断邮箱是否已存在？
	// 并发查询的时候多个发现不存在，后插入的出现错误❌
	// 加锁的话会有性能问题，邮箱冲突的情况很少
	return err
}

// User 直接对于数据库表结构，两者一一对应
// 如 entity modle PO (persistent object) ...
type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 唯一索引 允许有多个空值
	// 但是不能有多个 ""
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`
	// 创建时间和更新时间，毫秒数，UTC
	// 避免应用代码和数据库时区的不一致性
	Ctime int64
	Utime int64
}
