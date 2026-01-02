package model

import (
	"context"

	"gorm.io/gorm"
)

// MessageModel 是对 message 表操作的封装接口
type MessageModel interface {
	Insert(ctx context.Context, msg *Message) error
	GetMessagesByUser(ctx context.Context, userId int64) ([]*Message, error)
	// 可以继续扩展更多函数：UpdateStatus、DeleteMessage、GetUnreadCount 等
}

// messageModel 实现了 MessageModel 接口
type messageModel struct {
	db *gorm.DB
}

// NewMessageModel 创建一个 MessageModel 实例
func NewMessageModel(db *gorm.DB) MessageModel {
	return &messageModel{db: db}
}

// Insert 新增一条消息
func (m *messageModel) Insert(ctx context.Context, msg *Message) error {
	return m.db.WithContext(ctx).Create(msg).Error
}

// GetMessagesByUser 根据用户 ID 查询收到的消息
func (m *messageModel) GetMessagesByUser(ctx context.Context, userId int64) ([]*Message, error) {
	var msgs []*Message
	err := m.db.WithContext(ctx).Where("to_user_id = ?", userId).Find(&msgs).Error
	return msgs, err
}

// 查询两个用户之间的所有消息
func (m *messageModel) GetChatMessagesBetween(ctx context.Context, userId1, userId2 int64) ([]*Message, error) {
	var msgs []*Message
	err := m.db.WithContext(ctx).
		Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
			userId1, userId2, userId2, userId1).
		Order("created_at asc").
		Find(&msgs).Error
	return msgs, err
}
