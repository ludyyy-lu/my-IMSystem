package push

import (
	"encoding/json"
	"log"

	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/session"
)

type Service struct {
	dispatcher *Dispatcher
}

func NewService(manager *session.Manager) *Service {
	return &Service{dispatcher: NewDispatcher(manager)}
}

func (s *Service) PushToUser(userID int64, messageType string, payload interface{}) {
	if s == nil || s.dispatcher == nil {
		return
	}
	data, err := json.Marshal(model.PushMessage{Type: messageType, Payload: payload})
	if err != nil {
		log.Printf("failed to marshal push payload: %v", err)
		return
	}
	if err := s.dispatcher.DispatchToUser(userID, data); err != nil {
		log.Printf("failed to push to user %d: %v", userID, err)
	}
}

func (s *Service) Broadcast(messageType string, payload interface{}) {
	if s == nil || s.dispatcher == nil {
		return
	}
	data, err := json.Marshal(model.PushMessage{Type: messageType, Payload: payload})
	if err != nil {
		log.Printf("failed to marshal broadcast payload: %v", err)
		return
	}
	s.dispatcher.Broadcast(data)
}
