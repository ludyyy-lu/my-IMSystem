// ws/push.go
package ws

import (
	"encoding/json"
	"log"
	"my-IMSystem/ws-gateway/internal/conn"

	"github.com/gorilla/websocket"
)

func PushToUser(userId int64, payload interface{}) {
	type PushMessage struct {
		Type    string      `json:"type"`    // e.g. "friend_event"
		Payload interface{} `json:"payload"` // 原始事件体
	}
	pushData := PushMessage{
		Type:    "friend_event",
		Payload: payload,
	}
	// data, err := json.Marshal(pushData)

	// 将消息对象编码为 JSON 字节流
	data, err := json.Marshal(pushData)
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
