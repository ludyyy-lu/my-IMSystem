// ws/push.go
package ws

import (
	"my-IMSystem/ws-gateway/internal/conn"
	"encoding/json"
	"log"
	"github.com/gorilla/websocket"
)

func PushToUser(userId int64, payload interface{}) {
	// 将消息对象编码为 JSON 字节流
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal push payload: %v", err)
		return
	}
	// 尝试从连接管理器中获取用户连接
	c, ok := conn.GlobalConnManager.Get(userId)
	if !ok || c == nil {
		log.Printf("User %d not connected, skipping push", userId)
		return
	}

	// 向用户写入 WebSocket 消息
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("Failed to push message to user %d: %v", userId, err)
		// 💡 出错时，你可以选择从连接池中移除这个连接
		conn.GlobalConnManager.Remove(userId)
	}
}
