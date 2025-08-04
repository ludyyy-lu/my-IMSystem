package model

type User struct {
	ID       int64  `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"` // 加密存储
}
