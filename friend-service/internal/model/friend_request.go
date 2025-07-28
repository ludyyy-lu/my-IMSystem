package model

import "time"

type FriendRequest struct {
	ID          int64     `gorm:"primaryKey"`
	FromUserID  int64     `gorm:"index"`                      // 发起方
	ToUserID    int64     `gorm:"index"`                      // 接收方
	Remark      string    `gorm:"type:varchar(255)"`
	Status      string    `gorm:"type:varchar(20);default:'pending'"` // pending/accepted/rejected
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
