// Package session_test 包含对 session.Session 生命周期的完整测试。
//
// # WebSocket 链路流程说明
//
// ## 服务器→客户端推送流程（Server Push）
//
//  1. 上层调用 sess.Send(payload)，将数据放入带缓冲的 send channel（容量 256）。
//  2. writeLoop goroutine 从 send channel 取出数据。
//  3. writeLoop 调用 conn.WriteMessage(TextMessage, data)，向底层 TCP 写入 WebSocket 文本帧。
//  4. gorilla/websocket 完成帧格式化后交给操作系统 TCP 缓冲区。
//  5. 客户端调用 clientConn.ReadMessage()，从 TCP 读取帧并返回原始 payload。
//
// ## 客户端→服务器发送流程（Client Send）
//
//  1. 客户端调用 clientConn.WriteMessage(TextMessage, data)。
//  2. 服务器端 readLoop goroutine 调用 conn.ReadMessage() 收到数据。
//  3. readLoop 调用 onMessage(userID, data) 回调，将消息交给 router 层处理。
//
// ## 心跳机制（Ping/Pong）
//
//   - writeLoop 每 50s 发送一次 Ping 帧，保持连接活跃。
//   - 客户端（gorilla 默认行为）自动回复 Pong。
//   - PongHandler 在收到 Pong 后重置 ReadDeadline，延长 60s 超时窗口。
package session_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"my-IMSystem/ws-gateway/internal/session"

	"github.com/gorilla/websocket"
)

var testUpgrader = websocket.Upgrader{
	CheckOrigin: func(*http.Request) bool { return true },
}

// makeWSPair 创建一对测试用 WebSocket 连接（serverConn + clientConn）。
// 底层通过 httptest.Server 完成 HTTP→WS 升级，模拟真实网络链路。
// 测试结束时 httptest.Server 会被 t.Cleanup 自动关闭。
func makeWSPair(t *testing.T) (serverConn, clientConn *websocket.Conn) {
	t.Helper()
	serverConnCh := make(chan *websocket.Conn, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := testUpgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Errorf("WS upgrade error: %v", err)
			return
		}
		serverConnCh <- conn
	}))
	t.Cleanup(srv.Close)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/"
	cl, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("WS dial: %v", err)
	}
	return <-serverConnCh, cl
}

// TestSession_ServerSendsToClient 验证服务器→客户端的推送链路是否完整可用。
//
// 链路：sess.Send(data) → send channel → writeLoop → conn.WriteMessage → TCP帧 → clientConn.ReadMessage
func TestSession_ServerSendsToClient(t *testing.T) {
	serverConn, clientConn := makeWSPair(t)
	t.Cleanup(func() { clientConn.Close() })

	sess := session.NewSession(1, serverConn, nil, nil)
	sess.Start()
	t.Cleanup(sess.Close)

	want := []byte(`{"type":"chat_message","payload":"hello WS"}`)
	if err := sess.Send(want); err != nil {
		t.Fatalf("sess.Send: %v", err)
	}

	if err := clientConn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, got, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("client ReadMessage: %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("received %q, want %q", got, want)
	}
}

// TestSession_ClientSendsToServer 验证客户端→服务器的发送链路是否完整可用。
//
// 链路：clientConn.WriteMessage → TCP帧 → readLoop.ReadMessage → onMessage(uid, payload)
func TestSession_ClientSendsToServer(t *testing.T) {
	serverConn, clientConn := makeWSPair(t)
	t.Cleanup(func() { clientConn.Close() })

	received := make(chan []byte, 1)
	sess := session.NewSession(1, serverConn, func(_ int64, data []byte) {
		received <- data
	}, nil)
	sess.Start()
	t.Cleanup(sess.Close)

	want := []byte(`{"type":"chat","to":2,"content":"hello server"}`)
	if err := clientConn.WriteMessage(websocket.TextMessage, want); err != nil {
		t.Fatalf("client WriteMessage: %v", err)
	}

	select {
	case got := <-received:
		if string(got) != string(want) {
			t.Errorf("onMessage got %q, want %q", got, want)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout: onMessage was never called")
	}
}

// TestSession_MultipleMessages 验证多条消息按发送顺序经 WS 链路依次送达客户端。
func TestSession_MultipleMessages(t *testing.T) {
	serverConn, clientConn := makeWSPair(t)
	t.Cleanup(func() { clientConn.Close() })

	sess := session.NewSession(1, serverConn, nil, nil)
	sess.Start()
	t.Cleanup(sess.Close)

	msgs := []string{"first", "second", "third"}
	for _, m := range msgs {
		if err := sess.Send([]byte(m)); err != nil {
			t.Fatalf("Send %q: %v", m, err)
		}
	}

	if err := clientConn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	for _, want := range msgs {
		_, got, err := clientConn.ReadMessage()
		if err != nil {
			t.Fatalf("ReadMessage: %v", err)
		}
		if string(got) != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

// TestSession_Close_Idempotent 验证 Close 的幂等性：
// 无论串行还是并发调用多少次，onClose 回调只执行一次，连接只关闭一次。
func TestSession_Close_Idempotent(t *testing.T) {
	serverConn, clientConn := makeWSPair(t)
	t.Cleanup(func() { clientConn.Close() })

	var closeCount atomic.Int32
	sess := session.NewSession(42, serverConn, nil, func(_ int64) {
		closeCount.Add(1)
	})
	sess.Start()

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sess.Close()
		}()
	}
	wg.Wait()

	if n := closeCount.Load(); n != 1 {
		t.Errorf("onClose called %d times, want exactly 1", n)
	}
}

// TestSession_SendAfterClose 验证 Session 关闭后 Send 返回错误，不会阻塞或 panic。
func TestSession_SendAfterClose(t *testing.T) {
	serverConn, clientConn := makeWSPair(t)
	t.Cleanup(func() { clientConn.Close() })

	sess := session.NewSession(1, serverConn, nil, nil)
	sess.Start()
	sess.Close()

	// 等待 context 取消信号传播
	time.Sleep(20 * time.Millisecond)

	if err := sess.Send([]byte("should fail")); err == nil {
		t.Error("Send after Close should return an error, got nil")
	}
}
