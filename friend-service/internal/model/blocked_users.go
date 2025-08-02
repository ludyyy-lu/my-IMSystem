package model

import "time"

type BlockedUser struct {
	ID        int64 `gorm:"primaryKey"`
	UserID    int64 `gorm:"index"` // 拉黑者
	BlockedID int64 `gorm:"index"` // 被拉黑人
	CreatedAt time.Time
}
