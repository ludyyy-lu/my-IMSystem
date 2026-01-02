// internal/ws/upgrader.go

package ws

import (
	"context"
	"errors"
	"my-IMSystem/auth-service/auth"
	"net/http"

	"github.com/gorilla/websocket"
)

// 「连接升级」和「token 鉴权」

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 放开跨域限制，生产环境要收紧！
	},
}


func UpgradeAndAuthenticate(w http.ResponseWriter, r *http.Request, authRpc auth.AuthClient) (*websocket.Conn, int64, error) {
	token := r.Header.Get("Sec-Websocket-Protocol")
	if token == "" {
		return nil, 0, errors.New("missing token in Sec-Websocket-Protocol")
	}

	// 通过 authRpc 远程校验
	resp, err := authRpc.VerifyToken(context.Background(), &auth.VerifyTokenReq{
		AccessToken: token,
	})
	if err != nil || !resp.Valid {
		return nil, 0, errors.New("invalid token")
	}

	conn, err := upgrader.Upgrade(w, r, http.Header{
		"Sec-WebSocket-Protocol": []string{token},
	})
	if err != nil {
		return nil, 0, err
	}

	return conn, resp.UserId, nil
}
