package session

import (
	"errors"
	"sync"
)

type Manager struct {
	mu       sync.RWMutex
	sessions map[int64]*Session
}

func NewManager() *Manager {
	return &Manager{sessions: make(map[int64]*Session)}
}

func (m *Manager) Add(session *Session) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[session.UserID] = session
}

func (m *Manager) Remove(userID int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.sessions, userID)
}

func (m *Manager) Get(userID int64) (*Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[userID]
	return session, ok
}

func (m *Manager) SendTo(userID int64, data []byte) error {
	session, ok := m.Get(userID)
	if !ok {
		return errors.New("session not found")
	}
	return session.Send(data)
}

func (m *Manager) Broadcast(data []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, sess := range m.sessions {
		_ = sess.Send(data)
	}
}

// IsOnline reports whether the user has an active session registered with
// this Manager.  It reflects in-process state only; for cross-node presence
// awareness use a PresenceStore backed by Redis.
func (m *Manager) IsOnline(userID int64) bool {
	_, ok := m.Get(userID)
	return ok
}

// Count returns the number of currently registered sessions.
func (m *Manager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.sessions)
}
