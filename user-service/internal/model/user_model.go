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

// 这有啥用？
func (m *UserModel) UpdateUser(ctx context.Context, user *User) error {
	return m.db.WithContext(ctx).Save(user).Error
}

func (m *UserModel) UpdateByID(ctx context.Context, id int64, updateData map[string]interface{}) error {
	return m.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(updateData).Error
}

// FindByIDs 批量查询用户信息
// 这个方法可以用来查询多个用户的详细信息，避免多次查询数据库
func (m *UserModel) FindByIDs(ctx context.Context, ids []int64) ([]*User, error) {
	var users []*User
	if err := m.db.WithContext(ctx).Where("id IN ?", ids).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
