package model

import "time"

type Message struct {
	Id         int64     `gorm:"primaryKey" json:"id"`
	FromUserId int64     `gorm:"column:from_user_id" json:"from_user_id"`
	ToUserId   int64     `gorm:"column:to_user_id" json:"to_user_id"`
	Content    string    `gorm:"column:content" json:"content"`
	MsgType    int       `gorm:"column:msg_type" json:"msg_type"`
	Status     int       `gorm:"column:status" json:"status"` // 0 未读，1 已读
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Message) TableName() string {
	return "message"
}