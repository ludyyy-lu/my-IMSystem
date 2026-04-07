// Package consume reads messages from Kafka topics and forwards them to online
// users via the push service.  Each topic runs its own consumer group in a
// dedicated goroutine.
package consume

import (
	"context"
	"encoding/json"
	"log"

	"my-IMSystem/common/common_model"
	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/push"

	"github.com/segmentio/kafka-go"
)

// StartConsumers starts Kafka consumers for chat messages and friend events.
// Both topics are consumed concurrently.  The returned cancel func stops all
// consumers; it should be called (e.g. via defer) when the gateway shuts down.
func StartConsumers(brokers []string, chatTopic, friendTopic string, pushService *push.Service) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	if chatTopic != "" {
		go startConsumer(ctx, brokers, chatTopic, "ws-gateway-chat-group", func(value []byte) {
			var msg common_model.ChatMessage
			if err := json.Unmarshal(value, &msg); err != nil {
				log.Printf("failed to parse chat message: %v", err)
				return
			}
			pushService.PushToUser(msg.ToUserId, model.PushTypeChatMessage, msg)
		})
	}
	if friendTopic != "" {
		go startConsumer(ctx, brokers, friendTopic, "ws-gateway-friend-group", func(value []byte) {
			var event common_model.FriendEvent
			if err := json.Unmarshal(value, &event); err != nil {
				log.Printf("failed to parse friend event: %v", err)
				return
			}
			log.Printf("[FriendEvent] %+v", event)
			pushService.PushToUser(event.ToUser, model.PushTypeFriendEvent, event)
		})
	}
	return cancel
}

// startConsumer runs a blocking Kafka read loop in the calling goroutine.
// handler is invoked with the raw message value for each successfully read
// Kafka record.
func startConsumer(ctx context.Context, brokers []string, topic, groupID string, handler func([]byte)) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	defer reader.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m, err := reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("Kafka read error (topic=%s): %v", topic, err)
			continue
		}
		handler(m.Value)
	}
}
