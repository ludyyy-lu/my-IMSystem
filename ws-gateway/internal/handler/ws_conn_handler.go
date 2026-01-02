package handler

import (
	"net/http"
	"strings"

	"my-IMSystem/auth-service/auth"
	"my-IMSystem/ws-gateway/internal/logic"
	"my-IMSystem/ws-gateway/internal/svc"

	"github.com/gorilla/websocket"

	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域，后面可以做限制
	},
}

func connectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. 提取 token（从 header 或 query 参数）
		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			http.Error(w, "unauthorized: token is required", http.StatusUnauthorized)
			return
		}
		token = strings.TrimPrefix(token, "Bearer ")

		// 2. 远程调用 Auth 服务校验 token
		resp, err := svcCtx.AuthRpc.VerifyToken(r.Context(), &auth.VerifyTokenReq{
			AccessToken: token,
		})
		if err != nil || !resp.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		userId := resp.UserId

		// 3. 升级为 WebSocket 连接
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}
		logx.Infof("WebSocket connection established for user ID: %d", userId)
		// 4. 将连接和 userId 交给逻辑层处理
		go logic.HandleWebSocketConnection(svcCtx, userId, conn)
	}
}
