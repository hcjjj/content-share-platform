package dao

import "gorm.io/gorm"

type CodeDAO struct {
	db *gorm.DB
}

func NewCodeDAO(db *gorm.DB) *CodeDAO {
	return &CodeDAO{
		db: db,
	}
}
