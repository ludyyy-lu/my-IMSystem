// Package router dispatches inbound WebSocket messages to the appropriate
// business-logic handler based on message type.
// It sits between the transport layer (which handles raw WebSocket I/O) and
// the downstream services (Kafka, chat RPC).
package router

import (
	"context"
	"encoding/json"
	"time"

	"my-IMSystem/chat-service/chat"
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
	default:
		logx.Errorf("unknown message type from user %d: %s", userID, msg.Type)
	}
}

// handleChat forwards a chat message to Kafka for async processing by chat-service.
func handleChat(svcCtx *svc.ServiceContext, fromUserID int64, msg model.WsMessage) {
	msg.From = fromUserID
	if svcCtx.Config.Kafka.Topic == "" {
		logx.Error("Kafka chat topic is not configured")
		return
	}
	if err := kafka.SendMessage(svcCtx.Config.Kafka.Topic, msg); err != nil {
		logx.Errorf("failed to enqueue chat message to Kafka: %v", err)
	}
}

// handleAck acknowledges a delivered message via the chat RPC service.
func handleAck(svcCtx *svc.ServiceContext, fromUserID int64, msg model.WsMessage) {
	if svcCtx.ChatRpc == nil {
		logx.Error("chat RPC client not initialized")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := svcCtx.ChatRpc.AckMessage(ctx, &chat.AckMessageReq{
		MessageId: msg.Content,
		UserId:    fromUserID,
	}); err != nil {
		logx.Errorf("failed to ACK message %s from user %d: %v", msg.Content, fromUserID, err)
		return
	}
	logx.Infof("ACK sent for message %s from user %d", msg.Content, fromUserID)
}
