// Package transport_test 包含对 WebSocket 接入层（ConnectHandler）的集成测试。
//
// # 完整 WebSocket 链路流程（端到端说明）
//
// ## 一、连接建立（HTTP → WebSocket 升级）
//
//  1. 客户端发起 HTTP GET /ws/connect，在 Authorization 头（或 query token）中携带 Bearer Token。
//  2. ConnectHandler 提取 token，调用 AuthService.VerifyToken(ctx, token) 验证身份。
//  3. 验证通过后，upgrader.Upgrade(w, r, header) 完成 HTTP 101 Switching Protocols 升级。
//  4. 创建 session.Session（包含 read/write goroutine 和心跳定时器），注册到 SessionManager。
//  5. 异步加载并推送离线消息（若 OfflineStore != nil）。
//
// ## 二、服务器→客户端推送流程
//
//	Kafka 消费者接收到消息后：
//	  consume.handler(value)
//	    → push.Service.PushToUser(uid, msgType, payload)
//	        → json.Marshal(PushMessage{Type, Payload})   // 序列化
//	        → session.Manager.SendTo(uid, data)          // 按 uid 路由
//	            → session.Session.Send(data)             // 入 send channel（cap 256）
//	            → writeLoop → conn.WriteMessage          // 写 WS 文本帧到 TCP
//	            → 客户端 conn.ReadMessage                 // 客户端解码 WS 帧
//
// ## 三、客户端→服务器发送流程
//
//	客户端 conn.WriteMessage(TextMessage, jsonBytes)
//	  → readLoop: conn.ReadMessage()                     // 服务器读取原始 WS 帧字节
//	  → router.HandleMessage(svcCtx, uid, payload)
//	      → json.Unmarshal → WsMessage{Type, To, Content, From}
//	      → switch msg.Type:
//	          "chat" → kafka.SendMessage(topic, msg)     // 异步发布到 Kafka（topic="" 时跳过）
//	          "ack"  → ChatRpc.AckMessage(...)           // gRPC 确认（ChatRpc=nil 时跳过）
//
// ## 四、连接关闭流程
//
//	readLoop 读取到错误（EOF / 连接被关闭）
//	  → session.Close()
//	      → context.cancel()                            // 通知 writeLoop 退出
//	      → onClose(uid) → SessionManager.Remove(uid)  // 路由表清理
//	      → close(send channel) + conn.Close()
package transport_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"my-IMSystem/auth-service/auth"
	"my-IMSystem/ws-gateway/internal/config"
	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/push"
	"my-IMSystem/ws-gateway/internal/rpc"
	"my-IMSystem/ws-gateway/internal/session"
	"my-IMSystem/ws-gateway/internal/svc"
	"my-IMSystem/ws-gateway/internal/transport"

	"github.com/gorilla/websocket"
	"google.golang.org/grpc"
)

// mockAuthClient 实现 auth.AuthClient 接口，用于测试中绕过真实 gRPC 调用。
// VerifyToken 仅对 validToken 字段值返回有效响应。
type mockAuthClient struct {
	validToken string
	userID     int64
}

func (m *mockAuthClient) VerifyToken(_ context.Context, in *auth.VerifyTokenReq, _ ...grpc.CallOption) (*auth.VerifyTokenResp, error) {
	if in.AccessToken == m.validToken {
		return &auth.VerifyTokenResp{Valid: true, UserId: m.userID}, nil
	}
	return &auth.VerifyTokenResp{Valid: false}, nil
}

func (m *mockAuthClient) Register(_ context.Context, _ *auth.RegisterReq, _ ...grpc.CallOption) (*auth.RegisterResp, error) {
	return nil, nil
}
func (m *mockAuthClient) Login(_ context.Context, _ *auth.LoginReq, _ ...grpc.CallOption) (*auth.LoginResp, error) {
	return nil, nil
}
func (m *mockAuthClient) ParseToken(_ context.Context, _ *auth.ParseTokenReq, _ ...grpc.CallOption) (*auth.ParseTokenResp, error) {
	return nil, nil
}
func (m *mockAuthClient) RefreshToken(_ context.Context, _ *auth.RefreshTokenReq, _ ...grpc.CallOption) (*auth.RefreshTokenResp, error) {
	return nil, nil
}
func (m *mockAuthClient) ListSessions(_ context.Context, _ *auth.ListSessionsReq, _ ...grpc.CallOption) (*auth.ListSessionsResp, error) {
	return nil, nil
}
func (m *mockAuthClient) LogoutSession(_ context.Context, _ *auth.LogoutSessionReq, _ ...grpc.CallOption) (*auth.LogoutSessionResp, error) {
	return nil, nil
}
func (m *mockAuthClient) GenerateToken(_ context.Context, _ *auth.GenerateTokenReq, _ ...grpc.CallOption) (*auth.GenerateTokenResp, error) {
	return nil, nil
}

// newTestSvcCtx 创建用于测试的 ServiceContext：
//   - AuthService：使用 mockAuthClient，绕过真实 gRPC。
//   - SessionManager：真实实现，负责 session 路由。
//   - PushService：真实实现，依赖 SessionManager。
//   - OfflineStore：nil，ConnectHandler 会安全跳过离线消息推送。
//   - ChatRpc：nil，router.handleAck 会安全跳过。
//   - Kafka Topic 为空（config.Config 零值），router.handleChat 会安全跳过入队。
func newTestSvcCtx(validToken string, userID int64) *svc.ServiceContext {
	mgr := session.NewManager()
	return &svc.ServiceContext{
		Config:         config.Config{},
		SessionManager: mgr,
		OfflineStore:   nil,
		AuthService:    rpc.NewAuthService(&mockAuthClient{validToken: validToken, userID: userID}),
		PushService:    push.NewService(mgr),
		ChatRpc:        nil,
	}
}

// newHandlerServer 启动带有 /ws/connect 路由的测试 HTTP 服务器。
func newHandlerServer(t *testing.T, svcCtx *svc.ServiceContext) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("/ws/connect", transport.ConnectHandler(svcCtx))
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// dialWS 向测试服务器 /ws/connect 发起 WebSocket 连接，携带 Bearer token。
func dialWS(t *testing.T, srv *httptest.Server, token string) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/connect"
	header := http.Header{"Authorization": []string{"Bearer " + token}}
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		t.Fatalf("WS dial: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

// waitForSession 轮询 Manager 直到 userID 对应的 session 出现或超时。
// ConnectHandler 在 WS 升级后同步注册 session，通常极快完成，
// 但用轮询保证在极端调度延迟下测试依然稳定。
func waitForSession(t *testing.T, mgr *session.Manager, userID int64, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, ok := mgr.Get(userID); ok {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("session for user %d not found within %v", userID, timeout)
}

// TestConnectHandler_MissingToken 验证没有携带 token 时 ConnectHandler 拒绝升级（HTTP 401）。
func TestConnectHandler_MissingToken(t *testing.T) {
	svcCtx := newTestSvcCtx("valid-token", 1)
	srv := newHandlerServer(t, svcCtx)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/connect"
	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err == nil {
		t.Fatal("expected error for missing token, got nil")
	}
	if resp != nil && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("HTTP status = %d, want 401", resp.StatusCode)
	}
}

// TestConnectHandler_InvalidToken 验证携带无效 token 时 ConnectHandler 拒绝升级（HTTP 401）。
func TestConnectHandler_InvalidToken(t *testing.T) {
	svcCtx := newTestSvcCtx("valid-token", 1)
	srv := newHandlerServer(t, svcCtx)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/connect"
	header := http.Header{"Authorization": []string{"Bearer wrong-token"}}
	_, resp, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err == nil {
		t.Fatal("expected error for invalid token, got nil")
	}
	if resp != nil && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("HTTP status = %d, want 401", resp.StatusCode)
	}
}

// TestConnectHandler_ServerPushToClient 验证完整的服务器→客户端推送链路：
//
//  1. 客户端持有效 token 建立 WS 连接（HTTP 升级 + 认证 + session 注册）。
//  2. ConnectHandler 完成认证和 session 注册后，push.Service 可路由到该用户。
//  3. push.Service.PushToUser 将消息序列化为 PushMessage JSON，通过 session.writeLoop 发送到客户端。
//  4. 客户端正确接收并解析 PushMessage，type 字段与发送时一致。
func TestConnectHandler_ServerPushToClient(t *testing.T) {
	const (
		validToken = "token-user-42"
		userID     = int64(42)
	)
	svcCtx := newTestSvcCtx(validToken, userID)
	srv := newHandlerServer(t, svcCtx)

	// Step 1: 客户端建立 WS 连接
	clientConn := dialWS(t, srv, validToken)

	// Step 2: 等待 session 在 Manager 中注册完毕
	waitForSession(t, svcCtx.SessionManager, userID, 3*time.Second)

	// Step 3: 模拟 Kafka 消费者通过 push.Service 向用户推送消息
	chatPayload := map[string]string{"from": "user1", "content": "hello from server"}
	svcCtx.PushService.PushToUser(userID, model.PushTypeChatMessage, chatPayload)

	// Step 4: 客户端接收并验证消息
	if err := clientConn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		t.Fatalf("SetReadDeadline: %v", err)
	}
	_, raw, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("client ReadMessage: %v", err)
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

// TestConnectHandler_ClientToServer 验证完整的客户端→服务器发送链路：
//
//  1. 客户端连接并发送 JSON chat 消息。
//  2. readLoop 接收消息并调用 router.HandleMessage。
//  3. router 解析 type="chat"，因 Kafka Topic 为空跳过入队（安全 no-op）。
//  4. 消息处理完毕后 session 保持活跃，连接不中断。
func TestConnectHandler_ClientToServer(t *testing.T) {
	const (
		validToken = "token-user-55"
		userID     = int64(55)
	)
	svcCtx := newTestSvcCtx(validToken, userID)
	srv := newHandlerServer(t, svcCtx)

	clientConn := dialWS(t, srv, validToken)
	waitForSession(t, svcCtx.SessionManager, userID, 3*time.Second)

	// 发送 chat 消息（Kafka topic 为空，router 会安全跳过入队，不会 panic）
	msg := model.WsMessage{Type: "chat", To: 100, Content: "test content"}
	data, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("json.Marshal: %v", err)
	}
	if err := clientConn.WriteMessage(websocket.TextMessage, data); err != nil {
		t.Fatalf("WriteMessage: %v", err)
	}

	// 给服务器端 readLoop 处理时间，然后验证 session 仍然在线
	time.Sleep(50 * time.Millisecond)
	if _, ok := svcCtx.SessionManager.Get(userID); !ok {
		t.Error("session should still be active after routing a message")
	}
}

// TestConnectHandler_SessionCleanupOnDisconnect 验证客户端断开后 session 从 Manager 中移除。
//
// 关闭流程：
//
//	clientConn.Close()
//	  → readLoop: conn.ReadMessage() 返回错误（websocket: close）
//	  → session.Close()
//	      → onClose(uid) → SessionManager.Remove(uid)   // 路由表清理
func TestConnectHandler_SessionCleanupOnDisconnect(t *testing.T) {
	const (
		validToken = "token-user-77"
		userID     = int64(77)
	)
	svcCtx := newTestSvcCtx(validToken, userID)
	srv := newHandlerServer(t, svcCtx)

	clientConn := dialWS(t, srv, validToken)
	waitForSession(t, svcCtx.SessionManager, userID, 3*time.Second)

	// 主动关闭客户端连接
	clientConn.Close()

	// 轮询等待 readLoop 检测到关闭并触发 onClose 回调
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if _, ok := svcCtx.SessionManager.Get(userID); !ok {
			return // 成功：session 已从 Manager 中移除
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Error("session should have been removed from manager after client disconnected")
}

// TestConnectHandler_TokenViaQueryString 验证通过 URL query 参数传递 token 的连接方式。
//
// 部分浏览器 WebSocket 客户端无法设置自定义 Header，改用 ?token=xxx 携带令牌。
func TestConnectHandler_TokenViaQueryString(t *testing.T) {
	const (
		validToken = "qtoken-user-88"
		userID     = int64(88)
	)
	svcCtx := newTestSvcCtx(validToken, userID)
	srv := newHandlerServer(t, svcCtx)

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/connect?token=" + validToken
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("WS dial with query token: %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	waitForSession(t, svcCtx.SessionManager, userID, 3*time.Second)
}
