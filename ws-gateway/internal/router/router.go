// Package router dispatches inbound WebSocket messages to the appropriate
// business-logic handler based on message type.
package router

import (
"context"
"encoding/json"
"fmt"
"time"

"my-IMSystem/chat-service/chat"
"my-IMSystem/common/common_model"
"my-IMSystem/common/kafka"
"my-IMSystem/ws-gateway/internal/model"
"my-IMSystem/ws-gateway/internal/svc"

"github.com/zeromicro/go-zero/core/logx"
)

// HandleMessage parses payload and dispatches to the correct handler.
func HandleMessage(svcCtx *svc.ServiceContext, userID int64, payload []byte) {
var msg model.WsMessage
if err := json.Unmarshal(payload, &msg); err != nil {
logx.Errorf("invalid message from user %d: %v", userID, err)
return
}
switch msg.Type {
case "chat":
handleChat(svcCtx, userID, msg)
case "ack":
handleAck(svcCtx, userID, msg)
case "batch_ack":
handleBatchAck(svcCtx, userID, msg)
default:
logx.Errorf("unknown message type from user %d: %s", userID, msg.Type)
}
}

// handleChat forwards a chat message to Kafka for async processing by chat-service
// and for real-time push to the receiver via the ws-gateway consumer.
func handleChat(svcCtx *svc.ServiceContext, fromUserID int64, msg model.WsMessage) {
if svcCtx.Config.Kafka.Topic == "" {
logx.Error("Kafka chat topic is not configured")
return
}
// Build a properly-typed ChatMessage so both the chat-service consumer
// (saves to DB) and the ws-gateway push consumer (delivers in real-time)
// can parse from_user_id / to_user_id correctly.
chatMsg := common_model.ChatMessage{
MessageId:  fmt.Sprintf("%d_%d_%d", fromUserID, msg.To, time.Now().UnixNano()),
FromUserId: fromUserID,
ToUserId:   msg.To,
Content:    msg.Content,
Timestamp:  time.Now().Unix(),
}
if err := kafka.SendMessage(svcCtx.Config.Kafka.Topic, chatMsg); err != nil {
logx.Errorf("failed to enqueue chat message to Kafka: %v", err)
}
}

// handleAck acknowledges a single message via the chat RPC service, then pushes
// an ack_result confirmation back to the requesting user.
func handleAck(svcCtx *svc.ServiceContext, fromUserID int64, msg model.WsMessage) {
if svcCtx.ChatRpc == nil {
logx.Error("chat RPC client not initialized")
return
}
// MessageId field takes precedence; fall back to Content for backward compatibility.
messageID := msg.MessageId
if messageID == "" {
messageID = msg.Content
}
if messageID == "" {
logx.Errorf("ack from user %d: missing message_id", fromUserID)
return
}

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

resp, err := svcCtx.ChatRpc.AckMessage(ctx, &chat.AckMessageReq{
MessageId: messageID,
UserId:    fromUserID,
})

success := err == nil && resp != nil && resp.Status == "OK"
if err != nil {
logx.Errorf("failed to ACK message %s from user %d: %v", messageID, fromUserID, err)
} else {
logx.Infof("ACK sent for message %s from user %d", messageID, fromUserID)
}

// 将 ack 结果推送回客户端
if svcCtx.PushService != nil {
svcCtx.PushService.PushToUser(fromUserID, model.PushTypeAckResult, map[string]interface{}{
"message_id": messageID,
"success":    success,
})
}
}

// handleBatchAck marks all messages in a conversation as read, then pushes a
// batch_ack_result back to the requesting user.
func handleBatchAck(svcCtx *svc.ServiceContext, fromUserID int64, msg model.WsMessage) {
if svcCtx.ChatRpc == nil {
logx.Error("chat RPC client not initialized")
return
}
if msg.PeerId == 0 {
logx.Errorf("batch_ack from user %d: missing peer_id", fromUserID)
return
}

ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := svcCtx.ChatRpc.BatchAckMessages(ctx, &chat.BatchAckReq{
UserId: fromUserID,
PeerId: msg.PeerId,
})

success := err == nil && resp != nil && resp.Status == "OK"
var ackedCount int64
if resp != nil {
ackedCount = resp.AckedCount
}
if err != nil {
logx.Errorf("batch_ack failed user=%d peer=%d: %v", fromUserID, msg.PeerId, err)
} else {
logx.Infof("batch_ack done user=%d peer=%d acked=%d", fromUserID, msg.PeerId, ackedCount)
}

if svcCtx.PushService != nil {
svcCtx.PushService.PushToUser(fromUserID, model.PushTypeBatchAckResult, map[string]interface{}{
"peer_id":     msg.PeerId,
"success":     success,
"acked_count": ackedCount,
})
}
}
