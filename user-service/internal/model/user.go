package model

import (
	"time"
)

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique;not null;size:64" json:"username"`
	Password  string    `gorm:"not null;size:128" json:"-"` // 加密后存储，响应时不返回
	Nickname  string    `gorm:"size:64" json:"nickname"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
