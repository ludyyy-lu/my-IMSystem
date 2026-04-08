package model

// PushMessage 是服务端通过 WebSocket 推送给客户端的通用消息格式。
type PushMessage struct {
Type    string      `json:"type"`
Payload interface{} `json:"payload"`
}

const (
// PushTypeChatMessage 推送类型：新聊天消息
PushTypeChatMessage = "chat_message"
// PushTypeFriendEvent 推送类型：好友事件
PushTypeFriendEvent = "friend_event"
// PushTypeOfflineMessage 推送类型：离线缓存消息（连接恢复后批量下发）
PushTypeOfflineMessage = "offline_message"
// PushTypeAckResult 推送类型：单条消息已读回执结果（推送给执行 ack 的用户）
PushTypeAckResult = "ack_result"
// PushTypeBatchAckResult 推送类型：批量已读回执结果（推送给执行 batch_ack 的用户）
PushTypeBatchAckResult = "batch_ack_result"
// PushTypeReadReceipt 推送类型：已读回执通知（推送给消息原始发送方）
PushTypeReadReceipt = "read_receipt"
)
