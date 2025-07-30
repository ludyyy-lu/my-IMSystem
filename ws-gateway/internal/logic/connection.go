package logic

import (
	"encoding/json"
	"io"
	"log"
	"time"

	"my-IMSystem/ws-gateway/internal/kafka"
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
	// ✅ 1. 加入连接池
	svcCtx.ConnManager.Add(userId, conn)
	defer func() {
		svcCtx.ConnManager.Remove(userId)
		conn.Close()
		log.Printf("User %d disconnected", userId)
	}()
	// ✅ 1.1 拉取 Redis 中的离线消息
	offlineMsgs, err := svcCtx.OfflineStore.LoadAndDelete(userId)
	if err != nil {
		log.Printf("Failed to load offline messages for user %d: %v", userId, err)
	} else {
		for _, m := range offlineMsgs {
			resp := map[string]interface{}{
				"type":    "chat",
				"from":    m.From,
				"content": m.Content,
			}
			respBytes, _ := json.Marshal(resp)
			if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
				log.Printf("Failed to deliver offline message to user %d: %v", userId, err)
			}
		}
	}
	// ✅ 2. 设置心跳超时 & Pong handler
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// ✅ 3. 启动 ping 心跳协程
	go func() {
		ticker := time.NewTicker(pingPeriod)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Ping failed for user %d: %v", userId, err)
					return
				}
			}
		}
	}()

	// ✅ 4. 接收前端消息（如聊天消息、请求）
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close for user %d: %v", userId, err)
			} else if err == io.EOF {
				log.Printf("Client closed connection user %d", userId)
			}
			break
		}

		var message model.Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Printf("Invalid message from user %d: %v\nRaw: %s", userId, err, string(msg))
			continue
		}

		log.Printf("Parsed message from %d: %+v", userId, message)
		// ✅ 使用消息分发器
		RouteMessage(svcCtx, userId, message)
	}
}

func handleChatMessage(svcCtx *svc.ServiceContext, fromUserId int64, msg model.Message) {
	// ✅ 设置发送者 ID
	msg.From = fromUserId

	// ✅ 发送到 Kafka（即便用户在线，也不在这里推送）
	err := kafka.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
	}
}


// 临时处理chat信息
// func handleChatMessage(svcCtx *svc.ServiceContext, fromUserId int64, msg model.Message) {
// 	toConn, _ := svcCtx.ConnManager.Get(msg.To)
// 	if toConn == nil {
// 		log.Printf("User %d is offline. Cannot deliver message.\n", msg.To)
// 		// TODO: 存入离线消息
// 		// ✅ 存入 Redis 离线消息
// 		msg.From = fromUserId // 补充发送者字段
// 		if err := svcCtx.OfflineStore.Save(msg.To, msg); err != nil {
// 			log.Printf("Failed to store offline message for user %d: %v", msg.To, err)
// 		}
// 		return
// 	}

// 	// 构建返回消息
// 	resp := map[string]interface{}{
// 		"type":    "chat",
// 		"from":    fromUserId,
// 		"content": msg.Content,
// 	}
// 	respBytes, _ := json.Marshal(resp)
// 	err := toConn.WriteMessage(websocket.TextMessage, respBytes)
// 	if err != nil {
// 		log.Printf("Failed to send message to user %d: %v", msg.To, err)
// 	}
// }
