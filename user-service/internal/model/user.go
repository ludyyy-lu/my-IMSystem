package model

import "time"

type User struct {
	ID        int64     `gorm:"primaryKey;autoIncrement"`
	Nickname  string    `gorm:"type:varchar(64);not null"`
	Avatar    string    `gorm:"type:varchar(255);default:''"`
	Bio       string    `gorm:"type:varchar(255);default:''"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Disabled  bool      `gorm:"default:false"`
	Gender    int32     `gorm:"default:0"` // 0未知 1男 2女
}
