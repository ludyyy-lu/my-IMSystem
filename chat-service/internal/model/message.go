package model

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

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

// GetChatMessages 查询两个用户之间的所有消息
// 查询两人之间的消息（双向）
// 只查早于某个时间戳之前的 （分页用）
// 限制数量，按时间倒序排
func GetChatMessages(db *gorm.DB, userId, peerId int64, limit int, before time.Time) ([]Message, error) {
	var messages []Message

	err := db.
		Where(
			"(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
			userId, peerId, peerId, userId,
		).
		Where("created_at < ?", before).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}

// 修改消息为已读
func MarkMessageAsRead(db *gorm.DB, messageID string, userID int64) error {
	id, err := strconv.ParseInt(messageID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid message ID: %v", err)
	}

	// 条件必须匹配当前用户是接收者
	result := db.Model(&Message{}).
		Where("id = ? AND to_user_id = ?", id, userID).
		Update("status", 1)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no message updated")
	}
	return nil
}
