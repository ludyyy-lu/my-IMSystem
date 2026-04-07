// Package push delivers messages to connected WebSocket sessions.
// It wraps session.Manager to provide typed push operations used by the Kafka
// consumers and other server-side senders.  When a user is offline (no active
// session), the message is forwarded to the OfflineStore so it can be delivered
// once the user reconnects.
package push

import (
	"encoding/json"
	"log"

	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/session"
)

// Service delivers messages to online users via their active sessions.
// If a target user is offline and an OfflineStore is configured, the message
// is buffered for later delivery.
type Service struct {
	manager      *session.Manager
	offlineStore session.OfflineStore
}

// NewService creates a Service backed by the given session manager.
// offlineStore may be nil; in that case offline messages are only logged.
func NewService(manager *session.Manager, offlineStore session.OfflineStore) *Service {
	return &Service{manager: manager, offlineStore: offlineStore}
}

// PushToUser serialises payload as a typed PushMessage and sends it to the
// user's active session.  If delivery fails (e.g. user offline), the message
// is saved to the OfflineStore when one is configured.
func (s *Service) PushToUser(userID int64, messageType string, payload interface{}) {
	data, err := json.Marshal(model.PushMessage{Type: messageType, Payload: payload})
	if err != nil {
		log.Printf("push marshal error for user %d: %v", userID, err)
		return
	}
	if err := s.manager.SendTo(userID, data); err != nil {
		log.Printf("push delivery failed for user %d (type=%s): %v – buffering offline message", userID, messageType, err)
		if s.offlineStore != nil {
			if saveErr := s.offlineStore.Save(userID, data); saveErr != nil {
				log.Printf("failed to save offline message for user %d: %v", userID, saveErr)
			}
		}
	}
}

// Broadcast serialises payload as a typed PushMessage and sends it to all
// currently connected users.  Offline storage is not applied to broadcasts.
func (s *Service) Broadcast(messageType string, payload interface{}) {
	data, err := json.Marshal(model.PushMessage{Type: messageType, Payload: payload})
	if err != nil {
		log.Printf("broadcast marshal error: %v", err)
		return
	}
	s.manager.Broadcast(data)
}
