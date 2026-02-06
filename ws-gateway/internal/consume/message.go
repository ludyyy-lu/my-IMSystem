package consume

import (
	"context"
	"encoding/json"
	"log"

	"my-IMSystem/common/common_model"
	"my-IMSystem/ws-gateway/internal/push"

	"github.com/segmentio/kafka-go"
)

func StartConsumers(brokers []string, chatTopic string, friendTopic string, pushService *push.Service) {
	if chatTopic != "" {
		go startChatConsumer(brokers, chatTopic, pushService)
	}
	if friendTopic != "" {
		go startFriendConsumer(brokers, friendTopic, pushService)
	}
}

func startChatConsumer(brokers []string, topic string, pushService *push.Service) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "ws-gateway-chat-group",
	})
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

		if pushService != nil {
			pushService.PushToUser(msg.ToUserId, "chat_message", msg)
		}
	}
}

func startFriendConsumer(brokers []string, topic string, pushService *push.Service) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: "ws-gateway-friend-group",
	})
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
		log.Printf("[FriendEvent] %+v", event)
		if pushService != nil {
			pushService.PushToUser(event.ToUser, "friend_event", event)
		}
	}
}
