package kafka

import (
	"context"
	"encoding/json"
	"log"

	"my-IMSystem/common/common_model"
	"my-IMSystem/ws-gateway/internal/ws1"

	"github.com/segmentio/kafka-go"
)

func StartFriendConsumer(brokers []string, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "ws-gateway-friend-group",
	})
	go func() {
		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Kafka read error: %v", err)
				continue
			}

			var event common_model.FriendEvent
			if err := json.Unmarshal(m.Value, &event); err != nil {
				log.Printf("Failed to parse friend event: %v", err)
				continue
			}
			log.Printf("[FriendEvent] %+v\n", event)
			// 推送给接收者（ToUserID）
			ws1.PushToUser(event.ToUser, event)
		}
	}()
}
