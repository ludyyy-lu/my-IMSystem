package model

import (
	"context"
	"gorm.io/gorm"
)

// 数据库访问接口
type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db: db}
}

func (m *UserModel) FindByID(ctx context.Context, id int64) (*User, error) {
	var user User
	if err := m.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
