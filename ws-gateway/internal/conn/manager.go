package conn

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnManager struct {
	conns sync.Map // map[int64]*websocket.Conn
}

func NewConnManager() *ConnManager {
	return &ConnManager{}
}

func (m *ConnManager) Add(uid int64, conn *websocket.Conn) {
	m.conns.Store(uid, conn)
}

func (m *ConnManager) Remove(uid int64) {
	m.conns.Delete(uid)
}

func (m *ConnManager) Get(uid int64) (*websocket.Conn, bool) {
	val, ok := m.conns.Load(uid)
	if !ok {
		return nil, false
	}
	return val.(*websocket.Conn), true
}

func (m *ConnManager) SendTo(uid int64, msg []byte) error {
	conn, ok := m.Get(uid)
	if !ok {
		return nil // 用户不在线，不报错
	}
	return conn.WriteMessage(websocket.TextMessage, msg)
}
