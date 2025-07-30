package handler

import (
	"log"
	"my-IMSystem/chat-service/internal/model"
)

// ChatMessageHandler 是 Kafka 消费 chat 消息时触发的业务处理逻辑
func ChatMessageHandler(msg *model.Message) {
	log.Printf("[Kafka] 收到聊天消息: %+v", msg)

	// TODO: 👇这里添加业务逻辑，比如：
	// - 写入 MySQL 聊天记录
	// - 存离线消息 Redis
	// - 回执推送给发送者（通过 gRPC 或 WebSocket 网关）

	// 示例打印（你可以换成任何处理逻辑）
	log.Printf("处理消息 from %d to %d 内容：%s", msg.FromUserId, msg.ToUserId, msg.Content)
}
