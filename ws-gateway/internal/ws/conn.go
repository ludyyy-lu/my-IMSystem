// internal/ws/conn.go
package ws

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Conn struct {
	UserID int64
	wsConn *websocket.Conn
	send   chan []byte
	mu     sync.Mutex // 保护 send 通道
}

func NewConn(userID int64, conn *websocket.Conn) *Conn {
	return &Conn{
		UserID: userID,
		wsConn: conn,
		send:   make(chan []byte, 256),
	}
}

func (c *Conn) Start(onClose func()) {
	go c.readLoop(onClose) // 启动读循环
	go c.writeLoop()
}

// func (c *Conn) Send(msg []byte) {
// 	select {
// 	case c.send <- msg:
// 	default:
// 		log.Printf("Send buffer full for user %d, dropping message", c.UserID)
// 	}
// }

func (c *Conn) Send(message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	select {
	case c.send <- message:
		return nil
	default:
		return errors.New("send channel full or closed")
	}
}

func (c *Conn) readLoop(onClose func()) {
	defer func() {
		// 通知连接管理器：当前连接已断开
		onClose()

		// 关闭连接
		c.Close()
	}()
	for {
		_, message, err := c.wsConn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}
		// TODO: 处理接收到的 message
		log.Printf("recv from user %d: %s", c.UserID, string(message))
	}
}

func (c *Conn) writeLoop() {
	ticker := time.NewTicker(time.Second * 30)
	defer func() {
		ticker.Stop()
		c.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				_ = c.wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			_ = c.wsConn.WriteMessage(websocket.TextMessage, msg)
		case <-ticker.C:
			_ = c.wsConn.WriteMessage(websocket.PingMessage, []byte{})
		}
	}
}

func (c *Conn) Close() {
	// 尝试关闭 channel（避免重复关闭 panic）
	select {
	case <-c.send:
		// already closed
	default:
		close(c.send)
	}

	// 关闭 websocket
	_ = c.wsConn.Close()
}
