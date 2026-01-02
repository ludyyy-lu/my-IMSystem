package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"my-IMSystem/chat-service/internal/model"
)
// 类型：接收到消息后的处理函数
type MessageHandlerFunc func(msg *model.Message)

func StartChatMessageConsumer(brokers []string, topic string, groupID string, handler MessageHandlerFunc) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})

	log.Println("[Kafka] Chat consumer started.")

	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				log.Printf("[Kafka] ReadMessage error: %v", err)
				continue
			}

			var chatMsg model.Message
			if err := json.Unmarshal(m.Value, &chatMsg); err != nil {
				log.Printf("[Kafka] Failed to unmarshal: %v", err)
				continue
			}

			handler(&chatMsg)
		}
	}()
}
