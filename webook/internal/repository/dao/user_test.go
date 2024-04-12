package dao

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name string
		// mock 的数据库
		mock func(t *testing.T) *sql.DB
		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "插入成功",
			// 这里是 sqlmock 不是 用 gomock
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(123, 1)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").
					WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			user: User{
				Nickname: sql.NullString{
					String: "hcjjj",
					Valid:  true,
				},
			},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这边要求传入的是 sql 的正则表达式
				// 1062 是唯一索引冲突
				// Insert 的实现逻辑 1062 → ErrUserDuplicate
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx:  context.Background(),
			user: User{
				//Nickname: "Tom",
			},
			wantErr: ErrUserDuplicate,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				// 这边要求传入的是 sql 的正则表达式
				mock.ExpectExec("INSERT INTO .*").
					WillReturnError(errors.New("数据库错误"))
				return db
			},
			ctx:  context.Background(),
			user: User{
				//Nickname: "Tom",
			},
			wantErr: errors.New("数据库错误"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDB := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDB,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			dao := NewUserDAO(db)
			err = dao.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
