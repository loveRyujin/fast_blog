package model

import (
	"gorm.io/gorm"

	"github.com/onexstack_practice/fast_blog/internal/pkg/rid"
)

// AfterCreate 在创建数据库记录之后生成 postID.
func (m *Post) AfterCreate(tx *gorm.DB) error {
	m.PostID = rid.PostID.New(uint64(m.ID))

	return tx.Save(m).Error
}

// AfterCreate 在创建数据库记录之后生成 userID.
func (m *User) AfterCreate(tx *gorm.DB) error {
	m.UserID = rid.UserID.New(uint64(m.ID))

	return tx.Save(m).Error
}
