package common_model

// ReadReceiptEvent 是单条消息已读回执事件，由 chat-service 在标记消息已读后
// 发布到 Kafka read-receipt topic。ws-gateway 消费此事件后将回执实时推送
// 给消息原始发送方，告知其消息已被阅读。
type ReadReceiptEvent struct {
// MessageId 是被标记为已读的消息 ID
MessageId string `json:"message_id"`
// SenderId 是消息的原始发送方（需要接收回执通知的用户）
SenderId int64 `json:"sender_id"`
// ReaderId 是执行已读操作的用户（消息接收方）
ReaderId int64 `json:"reader_id"`
// ReadAt 是已读操作发生的 Unix 时间戳（毫秒）
ReadAt int64 `json:"read_at"`
}

// BatchReadReceiptEvent 是批量已读回执事件，由 chat-service 在
// BatchAckMessages 操作完成后发布，包含了整个会话被标记为已读的汇总信息。
type BatchReadReceiptEvent struct {
// SenderId 是消息的原始发送方（需要接收回执通知的用户）
SenderId int64 `json:"sender_id"`
// ReaderId 是执行已读操作的用户（消息接收方）
ReaderId int64 `json:"reader_id"`
// AckedCount 是本次批量操作实际标记为已读的消息数量
AckedCount int64 `json:"acked_count"`
// ReadAt 是批量已读操作发生的 Unix 时间戳（毫秒）
ReadAt int64 `json:"read_at"`
}
