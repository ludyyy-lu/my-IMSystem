// internal/ws/manager.go
package ws

import (
	"sync"
)

type ConnManager struct {
	conns map[int64]*Conn // 用户ID → 连接
	lock  sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		conns: make(map[int64]*Conn),
	}
}

func (m *ConnManager) AddConn(userID int64, conn *Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.conns[userID] = conn
}

func (m *ConnManager) RemoveConn(userID int64) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.conns, userID)
}

func (m *ConnManager) GetConn(userID int64) (*Conn, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	conn, ok := m.conns[userID]
	return conn, ok
}

func (m *ConnManager) Broadcast(message []byte) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, conn := range m.conns {
		conn.Send(message) // 后面 Conn 结构体中定义 Send 方法
	}
}
