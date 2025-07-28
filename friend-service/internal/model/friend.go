package model

import "time"

type Friend struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"index"`     // 自己
	FriendID  int64     `gorm:"index"`     // 好友
	CreatedAt time.Time
}
