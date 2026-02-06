package transport

import (
	"context"
	"net/http"
	"strings"
	"time"

	"my-IMSystem/chat-service/chat"
	"my-IMSystem/common/kafka"
	"my-IMSystem/ws-gateway/internal/model"
	"my-IMSystem/ws-gateway/internal/session"
	"my-IMSystem/ws-gateway/internal/svc"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域，后面可以做限制
	},
}

func ConnectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, protocolToken := extractToken(r)
		if token == "" {
			http.Error(w, "unauthorized: token is required", http.StatusUnauthorized)
			return
		}

		if svcCtx.AuthService == nil {
			http.Error(w, "auth service unavailable", http.StatusInternalServerError)
			return
		}

		userId, err := svcCtx.AuthService.VerifyToken(r.Context(), token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		responseHeader := http.Header{}
		if protocolToken != "" {
			responseHeader.Set("Sec-WebSocket-Protocol", protocolToken)
		}
		conn, err := upgrader.Upgrade(w, r, responseHeader)
		if err != nil {
			http.Error(w, "failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		onMessage := func(uid int64, payload []byte) {
			handleMessage(svcCtx, uid, payload)
		}
		onClose := func(uid int64) {
			logx.Infof("WebSocket connection closed for user ID: %d", uid)
		}

		sess := session.NewSession(userId, conn, svcCtx.SessionManager, svcCtx.OfflineStore, onMessage, onClose)
		sess.Start()
		logx.Infof("WebSocket connection established for user ID: %d", userId)
	}
}

func extractToken(r *http.Request) (string, string) {
	token := r.Header.Get("Authorization")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	protocolToken := r.Header.Get("Sec-WebSocket-Protocol")
	if token == "" {
		token = protocolToken
	}
	return strings.TrimPrefix(token, "Bearer "), protocolToken
}

func handleMessage(svcCtx *svc.ServiceContext, userId int64, payload []byte) {
	msg, err := ParseMessage(payload)
	if err != nil {
		logx.Errorf("invalid message from user %d: %v", userId, err)
		return
	}

	switch msg.Type {
	case "chat":
		handleChatMessage(svcCtx, userId, msg)
	case "ack":
		handleAckMessage(svcCtx, userId, msg)
	default:
		logx.Errorf("unknown message type: %s", msg.Type)
	}
}

func handleChatMessage(svcCtx *svc.ServiceContext, fromUserId int64, msg model.WsMessage) {
	msg.From = fromUserId
	if svcCtx.Config.Kafka.Topic == "" {
		logx.Error("Kafka topic is not configured")
		return
	}
	if err := kafka.SendMessage(svcCtx.Config.Kafka.Topic, msg); err != nil {
		logx.Errorf("failed to send message to Kafka: %v", err)
	}
}

func handleAckMessage(svcCtx *svc.ServiceContext, fromUserId int64, msg model.WsMessage) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if svcCtx.ChatRpc == nil {
		logx.Error("chat rpc client not initialized")
		return
	}

	resp, err := svcCtx.ChatRpc.AckMessage(ctx, &chat.AckMessageReq{
		MessageId: msg.Content,
		UserId:    fromUserId,
	})
	if err != nil {
		logx.Errorf("failed to send ACK to chat-service: %v", err)
		return
	}
	logx.Infof("ACK sent for message %s from user %d | resp: %+v", msg.Content, fromUserId, resp)
}
