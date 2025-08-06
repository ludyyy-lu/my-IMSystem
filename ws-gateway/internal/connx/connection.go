package connx

import (
	"context"
	"fmt"
	"log"
	"my-IMSystem/ws-gateway/internal/svc"
	"time"

	"github.com/gorilla/websocket"
)

// Connection 代表一个 WebSocket 连接
// 包含用户 ID、连接对象、服务上下文、发送通道等
// 以及一个上下文和取消函数用于管理连接的生命周期
// 该结构体用于处理 WebSocket 连接的读写操作、心跳检测
const (
	pongWait   = 60 * time.Second
	pingPeriod = 50 * time.Second
)

type Connection struct {
	UserId   int64
	Conn     *websocket.Conn
	SvcCtx   *svc.ServiceContext
	SendChan chan []byte
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewConnection(svcCtx *svc.ServiceContext, userId int64, conn *websocket.Conn) *Connection {
	ctx, cancel := context.WithCancel(context.Background())
	return &Connection{
		UserId:   userId,
		Conn:     conn,
		SvcCtx:   svcCtx,
		SendChan: make(chan []byte, 100),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (c *Connection) Start() {
	c.SvcCtx.ConnManager.Add(c.UserId, c.Conn)

	// 心跳处理
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	go c.readLoop()
	go c.pingLoop()
	go c.loadOfflineMessages()
}

func (c *Connection) Close() {
	c.cancel()
	c.SvcCtx.ConnManager.Remove(c.UserId)
	c.Conn.Close()
	log.Printf("User %d disconnected", c.UserId)
}

func (c *Connection) loadOfflineMessages() {
	messages, err := c.SvcCtx.RedisClient.LRange(c.ctx, fmt.Sprintf("offline:%d", c.UserId), int64(0), int64(-1)).Result()
	if err != nil {
		log.Printf("failed to load offline messages: %v", err)
		return
	}
	for _, msg := range messages {
		err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Printf("failed to send offline message: %v", err)
		}
	}
	// 清空离线消息
	c.SvcCtx.RedisClient.Del(c.ctx, fmt.Sprintf("offline:%d", c.UserId))
}

// readLoop handles incoming messages from the WebSocket connection.
func (c *Connection) readLoop() {
	defer c.Close()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("read error for user %d: %v", c.UserId, err)
			break
		}
		// Here you can handle the received message, e.g., dispatch to handler
		log.Printf("Received message from user %d: %s", c.UserId, string(message))
	}
}

// pingLoop sends periodic ping messages to keep the connection alive.
func (c *Connection) pingLoop() {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("ping error for user %d: %v", c.UserId, err)
				c.Close()
				return
			}
		}
	}
}

func (c *Connection) Context() context.Context {
	return c.ctx
}
