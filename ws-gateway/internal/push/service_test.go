// Package push_test 包含对 push.Service 的测试。
//
// push.Service 是 ws-gateway 推送层的顶级 API，位于 Kafka 消费者与 WebSocket session 之间：
//
//	Kafka消费者 / 其他服务
//	  → push.Service.PushToUser(uid, msgType, payload)
//	      → json.Marshal(PushMessage{Type: msgType, Payload: payload})  // 序列化为 JSON
//	      → session.Manager.SendTo(uid, data)                           // 按 uid 路由
//	          → session.Session.Send(data)                              // 入 send channel
//	          → writeLoop → conn.WriteMessage                           // 写 WS 帧
//	          → 客户端 ReadMessage                                        // 客户端收到
package push_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/push"
	"my-IMSystem/ws-gateway/internal/session"

	"github.com/gorilla/websocket"
)

var pushTestUpgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

// mockOfflineStore is a thread-safe in-memory OfflineStore for tests.
type mockOfflineStore struct {
	mu   sync.Mutex
	msgs map[int64][][]byte
}

func (m *mockOfflineStore) Save(userID int64, data []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.msgs == nil {
		m.msgs = make(map[int64][][]byte)
	}
	m.msgs[userID] = append(m.msgs[userID], data)
	return nil
}

func (m *mockOfflineStore) LoadAndDelete(userID int64) ([][]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	msgs := m.msgs[userID]
	delete(m.msgs, userID)
	return msgs, nil
}

func (m *mockOfflineStore) get(userID int64) [][]byte {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.msgs[userID]
}

// makePushWSPair 创建一对测试用 WebSocket 连接，供 push 包测试使用。
func makePushWSPair(t *testing.T) (serverConn, clientConn *websocket.Conn) {
	t.Helper()
	ch := make(chan *websocket.Conn, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := pushTestUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("upgrade: %v", err)
			return
		}
		ch <- conn
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/"
	cl, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	return <-ch, cl
}

// TestPushService_PushToUser 验证 PushToUser 将结构化推送消息序列化后完整送达在线用户。
//
// 断言：
//  1. 客户端能收到消息（WS 链路可通）。
//  2. 消息 JSON 中 type 字段与调用时传入的 messageType 一致。
//  3. payload 字段不为空（原始数据被保留）。
func TestPushService_PushToUser(t *testing.T) {
	mgr := session.NewManager()
	svc := push.NewService(mgr, nil)

	serverConn, clientConn := makePushWSPair(t)
	t.Cleanup(func() { clientConn.Close() })

	sess := session.NewSession(1, serverConn, nil, nil)
	sess.Start()
	t.Cleanup(sess.Close)
	mgr.Add(sess)

	svc.PushToUser(1, model.PushTypeChatMessage, map[string]string{"text": "hi"})

	if err := clientConn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, raw, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage: %v", err)
	}

	var pm model.PushMessage
	if err := json.Unmarshal(raw, &pm); err != nil {
		t.Fatalf("unmarshal PushMessage: %v (%s)", err, raw)
	}
	if pm.Type != model.PushTypeChatMessage {
		t.Errorf("type = %q, want %q", pm.Type, model.PushTypeChatMessage)
	}
	if pm.Payload == nil {
		t.Error("payload should not be nil")
	}
}

// TestPushService_PushToUser_Offline 验证对离线（无 session）的用户调用 PushToUser
// 不会 panic，仅记录日志后静默返回（无 OfflineStore 时）。
func TestPushService_PushToUser_Offline(t *testing.T) {
	mgr := session.NewManager()
	svc := push.NewService(mgr, nil)
	// 用户 999 没有在线 session，PushToUser 应静默失败，不 panic
	svc.PushToUser(999, model.PushTypeChatMessage, "some data")
}

// TestPushService_PushToUser_OfflineStore 验证当用户离线时，消息被正确保存到 OfflineStore。
func TestPushService_PushToUser_OfflineStore(t *testing.T) {
	mgr := session.NewManager()
	store := &mockOfflineStore{}
	svc := push.NewService(mgr, store)

	// 用户 888 离线，PushToUser 应将消息存入 OfflineStore
	svc.PushToUser(888, model.PushTypeChatMessage, map[string]string{"text": "offline msg"})

	saved := store.get(888)
	if len(saved) != 1 {
		t.Fatalf("expected 1 saved message, got %d", len(saved))
	}
	var pm model.PushMessage
	if err := json.Unmarshal(saved[0], &pm); err != nil {
		t.Fatalf("unmarshal saved message: %v", err)
	}
	if pm.Type != model.PushTypeChatMessage {
		t.Errorf("saved type = %q, want %q", pm.Type, model.PushTypeChatMessage)
	}
}

// TestPushService_Broadcast 验证 Broadcast 将同一条消息广播到所有在线用户，
// 且每个客户端收到的 type 字段均正确。
func TestPushService_Broadcast(t *testing.T) {
	mgr := session.NewManager()
	svc := push.NewService(mgr, nil)

	type entry struct {
		client *websocket.Conn
	}
	entries := make([]entry, 3)
	for i := 0; i < 3; i++ {
		sc, cc := makePushWSPair(t)
		t.Cleanup(func() { cc.Close() })
		sess := session.NewSession(int64(200+i), sc, nil, nil)
		sess.Start()
		t.Cleanup(sess.Close)
		mgr.Add(sess)
		entries[i] = entry{cc}
	}

	svc.Broadcast(model.PushTypeFriendEvent, map[string]string{"msg": "broadcast"})

	for i, e := range entries {
		if err := e.client.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
			t.Fatalf("entry %d SetReadDeadline: %v", i, err)
		}
		_, raw, err := e.client.ReadMessage()
		if err != nil {
			t.Fatalf("entry %d ReadMessage: %v", i, err)
		}
		var pm model.PushMessage
		if err := json.Unmarshal(raw, &pm); err != nil {
			t.Fatalf("entry %d unmarshal: %v", i, err)
		}
		if pm.Type != model.PushTypeFriendEvent {
			t.Errorf("entry %d type = %q, want %q", i, pm.Type, model.PushTypeFriendEvent)
		}
	}
}
