// ws/push.go
package ws1

import (
	"encoding/json"
	"log"
	"my-IMSystem/ws-gateway/internal/conn"
)

func PushToUser(userId int64, messageType string, payload interface{}) {
	type PushMessage struct {
		Type    string      `json:"type"`    // e.g. "friend_event"
		Payload interface{} `json:"payload"` // 原始事件体
	}
	pushData := PushMessage{
		Type:    messageType,
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
	if conn.GlobalConnManager == nil {
		log.Printf("connection manager not initialized, skipping push")
		return
	}
	if _, ok := conn.GlobalConnManager.Get(userId); !ok {
		log.Printf("User %d not connected, skipping push", userId)
		return
	}
	if err := conn.GlobalConnManager.SendTo(userId, data); err != nil {
		log.Printf("Failed to push message to user %d: %v", userId, err)
	}
}
