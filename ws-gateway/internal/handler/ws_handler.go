// internal/handler/ws_handler.go
package handler

import (
	"net/http"

	"my-IMSystem/ws-gateway/internal/svc"
)

func WsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return connectHandler(svcCtx)
}
