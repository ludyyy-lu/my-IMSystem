package logic

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"my-IMSystem/chat-service/chat"
	"my-IMSystem/common/kafka"
	"my-IMSystem/ws-gateway/internal/connx"
	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/svc"
	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

// 每一个前端 WebSocket 客户端成功连接后，服务端用来处理这条连接的主函数
func HandleWebSocketConnection(svcCtx *svc.ServiceContext, userId int64, conn *websocket.Conn) {
	c := connx.NewConnection(svcCtx, userId, conn)
	defer c.Close()
	c.Start()

	// 监听消息
	for {
		select {
		case <-c.Context().Done():
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Printf("Read failed: %v", err)
				return
			}
			var message model.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				log.Printf("Invalid message: %v", err)
				continue
			}
			RouteMessage(svcCtx, userId, message)
		}
	}
}

func handleChatMessage(svcCtx *svc.ServiceContext, fromUserId int64, msg model.Message) {
	// ✅ 设置发送者 ID
	msg.From = fromUserId

	// ✅ 发送到 Kafka（即便用户在线，也不在这里推送）
	err := kafka.SendMessage(svcCtx.Config.Kafka.Topic, msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
	}
}

func handleAckMessage(svcCtx *svc.ServiceContext, fromUserId int64, msg model.Message) {
	messageID := msg.Content
	// 调用 chat-service 的 RPC 接口，传递 ACK 请求
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := svcCtx.ChatRpc.AckMessage(ctx, &chat.AckMessageReq{
		MessageId: messageID,
		UserId:    fromUserId,
	})
	if err != nil {
		log.Printf("❌ Failed to send ACK to chat-service: %v", err)
	} else {
		log.Printf("✅ ACK sent for message %s from user %d | resp: %+v", messageID, fromUserId, resp)
	}
}
