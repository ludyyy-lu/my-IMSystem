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

func StartConsumers(brokers []string, chatTopic string, friendTopic string, pushService *push.Service) context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	if chatTopic != "" {
		go startChatConsumer(ctx, brokers, chatTopic, pushService)
	}
	if friendTopic != "" {
		go startFriendConsumer(ctx, brokers, friendTopic, pushService)
	}
	return cancel
}

func startChatConsumer(ctx context.Context, brokers []string, topic string, pushService *push.Service) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "ws-gateway-chat-group",
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
			log.Printf("Kafka read error (chat): %v", err)
			continue
		}

		var msg common_model.ChatMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("Failed to parse chat message: %v", err)
			continue
		}

		if pushService != nil {
			pushService.PushToUser(msg.ToUserId, model.PushTypeChatMessage, msg)
		}
	}
}

func startFriendConsumer(ctx context.Context, brokers []string, topic string, pushService *push.Service) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "ws-gateway-friend-group",
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
			log.Printf("Kafka read error: %v", err)
			continue
		}

		var event common_model.FriendEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("Failed to parse friend event: %v", err)
			continue
		}
		log.Printf("[FriendEvent] %+v", event)
		if pushService != nil {
			pushService.PushToUser(event.ToUser, model.PushTypeFriendEvent, event)
		}
	}
}
