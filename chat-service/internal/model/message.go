package model

import (
"fmt"
"strconv"
"time"

"gorm.io/gorm"
)

// Message 是消息表的 GORM 模型。status: 0=未读，1=已读。
type Message struct {
Id         int64     `gorm:"primaryKey" json:"id"`
FromUserId int64     `gorm:"column:from_user_id;index:idx_from_to" json:"from_user_id"`
ToUserId   int64     `gorm:"column:to_user_id;index:idx_to_status" json:"to_user_id"`
Content    string    `gorm:"column:content" json:"content"`
MsgType    int       `gorm:"column:msg_type" json:"msg_type"`
Status     int       `gorm:"column:status;index:idx_to_status" json:"status"` // 0 未读，1 已读
CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
}

func (Message) TableName() string {
return "message"
}

// GetChatMessages 查询两个用户之间的消息（双向），只返回早于 before 的记录，
// 按时间倒序排列，最多返回 limit 条（分页用）。
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
return messages, err
}

// MarkMessageAsRead 将单条消息标记为已读（幂等：已是已读状态时不报错）。
// 仅允许消息的接收方（to_user_id）执行此操作。
func MarkMessageAsRead(db *gorm.DB, messageID string, userID int64) error {
id, err := strconv.ParseInt(messageID, 10, 64)
if err != nil {
return fmt.Errorf("invalid message ID: %w", err)
}
result := db.Model(&Message{}).
Where("id = ? AND to_user_id = ? AND status = 0", id, userID).
Update("status", 1)
return result.Error
}

// MarkMessageAsReadAndGet 将单条消息标记为已读，并返回消息对象（用于获取发送方 ID）。
// 若消息不存在或已是已读状态，返回 (nil, nil)，幂等安全。
func MarkMessageAsReadAndGet(db *gorm.DB, messageID string, userID int64) (*Message, error) {
id, err := strconv.ParseInt(messageID, 10, 64)
if err != nil {
return nil, fmt.Errorf("invalid message ID: %w", err)
}

var msg Message
if err := db.Where("id = ? AND to_user_id = ?", id, userID).First(&msg).Error; err != nil {
if err == gorm.ErrRecordNotFound {
return nil, nil
}
return nil, err
}
if msg.Status == 1 {
return nil, nil // 已读，幂等
}
if err := db.Model(&msg).Update("status", 1).Error; err != nil {
return nil, err
}
return &msg, nil
}

// BatchMarkConversationAsRead 将 peerId 发给 userId 的所有未读消息标记为已读，
// 返回实际更新的行数。幂等：重复调用安全。
func BatchMarkConversationAsRead(db *gorm.DB, userId, peerId int64) (int64, error) {
result := db.Model(&Message{}).
Where("from_user_id = ? AND to_user_id = ? AND status = 0", peerId, userId).
Update("status", 1)
return result.RowsAffected, result.Error
}

// UnreadConversationStat 是单个会话的未读统计快照。
type UnreadConversationStat struct {
PeerId        int64
UnreadCount   int64
LastContent   string
LastTimestamp time.Time
}

// GetUnreadCountByPeer 查询 userId 收到的所有未读消息，按发送方聚合统计，
// 同时返回每个会话最近一条未读消息的内容与时间。结果按最新消息时间降序。
func GetUnreadCountByPeer(db *gorm.DB, userId int64) ([]UnreadConversationStat, error) {
type row struct {
FromUserId int64
Cnt        int64
LastMsg    string
LastTime   time.Time
}

var rows []row
err := db.Model(&Message{}).
Select("from_user_id, COUNT(*) AS cnt, " +
"SUBSTRING_INDEX(GROUP_CONCAT(content ORDER BY created_at DESC SEPARATOR '\x00'), '\x00', 1) AS last_msg, " +
"MAX(created_at) AS last_time").
Where("to_user_id = ? AND status = 0", userId).
Group("from_user_id").
Order("last_time DESC").
Scan(&rows).Error
if err != nil {
return nil, err
}

stats := make([]UnreadConversationStat, 0, len(rows))
for _, r := range rows {
stats = append(stats, UnreadConversationStat{
PeerId:        r.FromUserId,
UnreadCount:   r.Cnt,
LastContent:   r.LastMsg,
LastTimestamp: r.LastTime,
})
}
return stats, nil
}
