package logic

import (
	"io"
	"log"
	"time"

	"my-IMSystem/ws-gateway/internal/svc"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

func HandleWebSocketConnection(svcCtx *svc.ServiceContext, userId int64, conn *websocket.Conn) {
	// ✅ 1. 加入连接池
	svcCtx.ConnManager.Add(userId, conn)
	defer func() {
		svcCtx.ConnManager.Remove(userId)
		conn.Close()
		log.Printf("User %d disconnected", userId)
	}()

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

		log.Printf("Received from %d: %s", userId, string(msg))

		// ✅ 示例处理逻辑（你未来要解析 msg 并路由到 chat-service 或 Kafka）
		// processMessage(userId, msg)
	}
}
