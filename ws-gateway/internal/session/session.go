package session

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"my-IMSystem/ws-gateway/internal/model"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

type Session struct {
	UserID       int64
	Conn         *websocket.Conn
	send         chan []byte
	ctx          context.Context
	cancel       context.CancelFunc
	manager      *Manager
	offlineStore *RedisOfflineMsgStore
	onMessage    func(int64, []byte)
	onClose      func(int64)
	state        State
	closeOnce    sync.Once
}

func NewSession(userID int64, conn *websocket.Conn, manager *Manager, offlineStore *RedisOfflineMsgStore, onMessage func(int64, []byte), onClose func(int64)) *Session {
	ctx, cancel := context.WithCancel(context.Background())
	return &Session{
		UserID:       userID,
		Conn:         conn,
		send:         make(chan []byte, 256),
		ctx:          ctx,
		cancel:       cancel,
		manager:      manager,
		offlineStore: offlineStore,
		onMessage:    onMessage,
		onClose:      onClose,
		state:        StateOnline,
	}
}

func (s *Session) Start() {
	if s.manager != nil {
		s.manager.Add(s)
	}

	s.Conn.SetReadDeadline(time.Now().Add(pongWait))
	s.Conn.SetPongHandler(func(string) error {
		s.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go s.readLoop()
	go s.writeLoop()
	go s.loadOfflineMessages()
}

func (s *Session) Send(message []byte) error {
	select {
	case <-s.ctx.Done():
		return errors.New("session closed")
	default:
	}

	select {
	case s.send <- message:
		return nil
	default:
		return errors.New("send channel full")
	}
}

func (s *Session) Close() {
	s.closeOnce.Do(func() {
		s.state = StateOffline
		s.cancel()
		if s.manager != nil {
			s.manager.Remove(s.UserID)
		}
		if s.onClose != nil {
			s.onClose(s.UserID)
		}
		close(s.send)
		_ = s.Conn.Close()
	})
}

func (s *Session) readLoop() {
	defer s.Close()
	for {
		_, message, err := s.Conn.ReadMessage()
		if err != nil {
			log.Printf("read error for user %d: %v", s.UserID, err)
			return
		}
		if s.onMessage != nil {
			s.onMessage(s.UserID, message)
		}
	}
}

func (s *Session) writeLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		s.Close()
	}()
	for {
		select {
		case <-s.ctx.Done():
			return
		case message, ok := <-s.send:
			if !ok {
				_ = s.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = s.Conn.WriteMessage(websocket.TextMessage, message)
		case <-ticker.C:
			if err := s.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ping error for user %d: %v", s.UserID, err)
				return
			}
		}
	}
}

func (s *Session) loadOfflineMessages() {
	if s.offlineStore == nil {
		return
	}
	messages, err := s.offlineStore.LoadAndDelete(s.UserID)
	if err != nil {
		log.Printf("failed to load offline messages for user %d: %v", s.UserID, err)
		return
	}
	for _, msg := range messages {
		pushMsg := model.PushMessage{Type: "offline_message", Payload: msg}
		payload, err := json.Marshal(pushMsg)
		if err != nil {
			log.Printf("failed to marshal offline message for user %d: %v", s.UserID, err)
			continue
		}
		_ = s.Send(payload)
	}
}
