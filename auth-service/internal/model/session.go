package model

import "time"

type Session struct {
	ID           int64     `gorm:"primaryKey"`
	Uid          int64     `gorm:"index"` // 用户ID
	DeviceId     string    `gorm:"index"` // 设备ID
	AccessToken  string    // 存储 AccessToken（可选）
	RefreshToken string    // 存储 RefreshToken（可选）
	ExpiresAt    time.Time // token 到期时间
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
