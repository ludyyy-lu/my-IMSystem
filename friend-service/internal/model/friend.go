package model

import "time"
// 好友关系表
// 记录用户之间的好友关系
type Friend struct {
	ID        int64     `gorm:"primaryKey"`
	UserID    int64     `gorm:"index"`     // 自己
	FriendID  int64     `gorm:"index"`     // 好友
	CreatedAt time.Time
}
