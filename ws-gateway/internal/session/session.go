// Package session manages the lifecycle of a single WebSocket connection,
// including heartbeat (ping/pong), inbound read loop, and outbound write loop.
// It has no knowledge of session registration, offline storage, or business logic;
// those concerns are handled by callers via callbacks.
package session

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

// Session represents a single authenticated WebSocket connection.
type Session struct {
	UserID    int64
	Conn      *websocket.Conn
	send      chan []byte
	ctx       context.Context
	cancel    context.CancelFunc
	onMessage func(int64, []byte)
	onClose   func(int64)
	state     State
	closeOnce sync.Once
}

// NewSession creates a new Session for the given user and connection.
// onMessage is called for each inbound message; onClose is called when the
// connection is closed (exactly once).
func NewSession(userID int64, conn *websocket.Conn, onMessage func(int64, []byte), onClose func(int64)) *Session {
	ctx, cancel := context.WithCancel(context.Background())
	return &Session{
		UserID:    userID,
		Conn:      conn,
		send:      make(chan []byte, 256),
		ctx:       ctx,
		cancel:    cancel,
		onMessage: onMessage,
		onClose:   onClose,
		state:     StateOnline,
	}
}

// Start begins the read and write goroutines for the session.
// The caller is responsible for registering the session with a Manager and
// delivering any pending offline messages before or after calling Start.
func (s *Session) Start() {
	if err := s.Conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Printf("failed to set read deadline for user %d: %v", s.UserID, err)
	}
	s.Conn.SetPongHandler(func(string) error {
		return s.Conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	go s.readLoop()
	go s.writeLoop()
}

// Send enqueues a message for delivery to the client.
// Returns an error if the session is closed or the send buffer is full.
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

// Close terminates the session (idempotent).
func (s *Session) Close() {
	s.closeOnce.Do(func() {
		s.state = StateOffline
		s.cancel()
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
