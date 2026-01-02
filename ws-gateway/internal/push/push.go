// internal/push/push.go
package push

import (
	"encoding/json"
	"log"
	"my-IMSystem/ws-gateway/internal/ws"
)

// PushToUser 推送二进制消息给某个在线用户
func PushToUser(userID int64, message []byte, connMgr *ws.ConnManager) {
	conn, ok := connMgr.GetConn(userID)
	if !ok {
		log.Printf("user %d not online, skip push\n", userID)
		return
	}

	err := conn.Send(message)
	if err != nil {
		log.Printf("failed to send message to user %d: %v\n", userID, err)
	}
}

// PushJSONToUser 将结构体编码为 JSON 并推送给用户
func PushJSONToUser(userID int64, data any, connMgr *ws.ConnManager) {
	msg, err := json.Marshal(data)
	if err != nil {
		log.Printf("[push] failed to marshal message for user %d: %v\n", userID, err)
		return
	}
	PushToUser(userID, msg, connMgr)
}
