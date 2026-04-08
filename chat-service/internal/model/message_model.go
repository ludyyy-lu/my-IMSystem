package model

import (
"context"

"gorm.io/gorm"
)

// MessageModel 是对 message 表操作的封装接口。
type MessageModel interface {
Insert(ctx context.Context, msg *Message) error
GetMessagesByUser(ctx context.Context, userId int64) ([]*Message, error)
GetChatMessagesBetween(ctx context.Context, userId1, userId2 int64) ([]*Message, error)
// GetUnreadStats 查询 userId 收到的所有未读消息，按会话对端聚合统计。
GetUnreadStats(ctx context.Context, userId int64) ([]UnreadConversationStat, error)
// BatchMarkRead 将 peerId 发给 userId 的所有未读消息标记为已读，返回实际更新条数。
BatchMarkRead(ctx context.Context, userId, peerId int64) (int64, error)
}

type messageModel struct {
db *gorm.DB
}

func NewMessageModel(db *gorm.DB) MessageModel {
return &messageModel{db: db}
}

func (m *messageModel) Insert(ctx context.Context, msg *Message) error {
return m.db.WithContext(ctx).Create(msg).Error
}

func (m *messageModel) GetMessagesByUser(ctx context.Context, userId int64) ([]*Message, error) {
var msgs []*Message
err := m.db.WithContext(ctx).Where("to_user_id = ?", userId).Find(&msgs).Error
return msgs, err
}

func (m *messageModel) GetChatMessagesBetween(ctx context.Context, userId1, userId2 int64) ([]*Message, error) {
var msgs []*Message
err := m.db.WithContext(ctx).
Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?)",
userId1, userId2, userId2, userId1).
Order("created_at asc").
Find(&msgs).Error
return msgs, err
}

func (m *messageModel) GetUnreadStats(ctx context.Context, userId int64) ([]UnreadConversationStat, error) {
return GetUnreadCountByPeer(m.db.WithContext(ctx), userId)
}

func (m *messageModel) BatchMarkRead(ctx context.Context, userId, peerId int64) (int64, error) {
return BatchMarkConversationAsRead(m.db.WithContext(ctx), userId, peerId)
}
