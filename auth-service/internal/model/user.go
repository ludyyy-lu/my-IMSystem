package model

type User struct {
	ID       int64  `gorm:"primaryKey"`
	Username string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Password string `gorm:"type:varchar(255);not null"` // 加密存储
	Nickname string `gorm:"type:varchar(64);not null"`
}
