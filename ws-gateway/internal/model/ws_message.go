package model

// WsMessage 是客户端通过 WebSocket 发送给服务端的消息格式。
type WsMessage struct {
// Type 是消息类型："chat" | "ack" | "batch_ack"
Type string `json:"type"`
// To 是接收方 userId（chat 消息时使用）
To int64 `json:"to,omitempty"`
// Content 是消息内容（chat 类型时使用）
Content string `json:"content,omitempty"`
// From 是发送方 userId（由服务端补充）
From int64 `json:"from,omitempty"`
// MessageId 是要 ACK 的消息 ID（ack 类型时使用）
MessageId string `json:"message_id,omitempty"`
// PeerId 是会话对端 ID（batch_ack 类型时使用）
PeerId int64 `json:"peer_id,omitempty"`
}
