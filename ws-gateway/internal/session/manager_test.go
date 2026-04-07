// manager_test.go 测试 session.Manager 的路由与广播能力。
//
// Manager 在整个 WS 链路中充当路由表：
//   - 用户上线时 ConnectHandler 调用 Add(sess) 注册 session。
//   - push.Service 或 Kafka 消费者通过 SendTo(uid, data) 将消息路由到指定用户。
//   - Broadcast(data) 将消息广播到所有在线用户。
//   - 用户断开时 onClose 回调调用 Remove(uid) 清理路由。
package session_test

import (
	"testing"
	"time"

	"my-IMSystem/ws-gateway/internal/session"

	"github.com/gorilla/websocket"
)

// newTestSession 创建一个测试用 Session（可选择是否启动 read/write goroutine）。
// 使用 makeWSPair 建立真实的 WebSocket 回环连接以满足 gorilla 的接口要求。
func newTestSession(t *testing.T, userID int64, start bool) (*session.Session, *websocket.Conn) {
	t.Helper()
	serverConn, clientConn := makeWSPair(t)
	t.Cleanup(func() { clientConn.Close() })
	sess := session.NewSession(userID, serverConn, nil, nil)
	if start {
		sess.Start()
		t.Cleanup(sess.Close)
	}
	return sess, clientConn
}

// TestManager_AddAndGet 验证 Manager 能正确存储并按 userID 检索 Session。
func TestManager_AddAndGet(t *testing.T) {
	m := session.NewManager()
	sess, _ := newTestSession(t, 10, false)
	m.Add(sess)

	got, ok := m.Get(10)
	if !ok {
		t.Fatal("Get(10): expected session to exist, got false")
	}
	if got != sess {
		t.Error("Get returned a different session pointer")
	}
}

// TestManager_Remove 验证 Remove 后 Get 返回 (nil, false)，路由条目已清除。
func TestManager_Remove(t *testing.T) {
	m := session.NewManager()
	sess, _ := newTestSession(t, 20, false)
	m.Add(sess)
	m.Remove(20)

	if _, ok := m.Get(20); ok {
		t.Error("Get after Remove should return false, but session still exists")
	}
}

// TestManager_SendTo_Online 验证 SendTo 将数据路由到在线用户的 Session，
// 并由 Session 的 writeLoop 最终写到对应的 WebSocket 客户端连接。
//
// 链路：Manager.SendTo(uid, data) → session.Send → send channel → writeLoop → clientConn.ReadMessage
func TestManager_SendTo_Online(t *testing.T) {
	m := session.NewManager()
	sess, clientConn := newTestSession(t, 30, true)
	m.Add(sess)

	payload := []byte(`{"type":"test","payload":"manager_send"}`)
	if err := m.SendTo(30, payload); err != nil {
		t.Fatalf("SendTo: %v", err)
	}

	if err := clientConn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, got, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("client ReadMessage: %v", err)
	}
	if string(got) != string(payload) {
		t.Errorf("got %q, want %q", got, payload)
	}
}

// TestManager_SendTo_Offline 验证 SendTo 对不存在（离线）的 userID 返回非 nil 错误。
func TestManager_SendTo_Offline(t *testing.T) {
	m := session.NewManager()
	if err := m.SendTo(999, []byte("data")); err == nil {
		t.Error("SendTo offline user should return an error, got nil")
	}
}

// TestManager_Broadcast 验证 Broadcast 将同一条消息发送到所有在线 Session。
//
// 链路：Manager.Broadcast(data) → 遍历所有 session → session.Send → writeLoop → 各客户端 ReadMessage
func TestManager_Broadcast(t *testing.T) {
	m := session.NewManager()

	type clientEntry struct {
		conn *websocket.Conn
	}
	clients := make([]clientEntry, 3)
	for i := 0; i < 3; i++ {
		sess, cl := newTestSession(t, int64(100+i), true)
		m.Add(sess)
		clients[i] = clientEntry{cl}
	}

	payload := []byte(`{"type":"broadcast","payload":"hello all"}`)
	m.Broadcast(payload)

	for i, c := range clients {
		if err := c.conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
			t.Fatalf("client %d SetReadDeadline: %v", i, err)
		}
		_, got, err := c.conn.ReadMessage()
		if err != nil {
			t.Fatalf("client %d ReadMessage: %v", i, err)
		}
		if string(got) != string(payload) {
			t.Errorf("client %d: got %q, want %q", i, got, payload)
		}
	}
}
