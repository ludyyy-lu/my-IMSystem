package kafka

import (
	"context"
	"encoding/json"
	"log"

	"my-IMSystem/common/common_model"
	"my-IMSystem/ws-gateway/internal/ws1"

	"github.com/segmentio/kafka-go"
)

func StartChatConsumer(brokers []string, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "ws-gateway-chat-group",
	})
	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Kafka read error (chat): %v", err)
				continue
			}

			var msg common_model.ChatMessage
			if err := json.Unmarshal(m.Value, &msg); err != nil {
				log.Printf("Failed to parse chat message: %v", err)
				continue
			}

			// 调用 WebSocket 推送
			ws1.PushToUser(msg.ToUserId, msg)
		}
	}()
}
