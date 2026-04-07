// Package push delivers messages to connected WebSocket sessions.
// It wraps session.Manager to provide typed push operations used by the Kafka
// consumers and other server-side senders.
package push

import (
	"encoding/json"
	"log"

	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/session"
)

// Service delivers messages to online users via their active sessions.
type Service struct {
	manager *session.Manager
}

// NewService creates a Service backed by the given session manager.
func NewService(manager *session.Manager) *Service {
	return &Service{manager: manager}
}

// PushToUser serialises payload as a typed PushMessage and sends it to the
// user's active session.  Errors (e.g. user offline) are logged and discarded.
func (s *Service) PushToUser(userID int64, messageType string, payload interface{}) {
	data, err := json.Marshal(model.PushMessage{Type: messageType, Payload: payload})
	if err != nil {
		log.Printf("push marshal error for user %d: %v", userID, err)
		return
	}
	if err := s.manager.SendTo(userID, data); err != nil {
		log.Printf("push delivery failed for user %d (type=%s): %v", userID, messageType, err)
	}
}

// Broadcast serialises payload as a typed PushMessage and sends it to all
// currently connected users.
func (s *Service) Broadcast(messageType string, payload interface{}) {
	data, err := json.Marshal(model.PushMessage{Type: messageType, Payload: payload})
	if err != nil {
		log.Printf("broadcast marshal error: %v", err)
		return
	}
	s.manager.Broadcast(data)
}
