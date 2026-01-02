// internal/handler/ws_handler.go
package handler

import (
	"log"
	"my-IMSystem/ws-gateway/internal/svc"
	"my-IMSystem/ws-gateway/internal/ws"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func WsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, userID, err := ws.UpgradeAndAuthenticate(w, r, svcCtx.AuthRpc)
		if err != nil {
			httpx.Error(w, err)
			return
		}

		c := ws.NewConn(userID, conn)
		svcCtx.ConnMgr.AddConn(userID, c)
		// 传入 onClose 函数，断开连接时从 ConnMgr 删除连接
		c.Start(func() {
			svcCtx.ConnMgr.RemoveConn(userID)
			log.Printf("connection closed for user %d", userID)
		})
	}
}
