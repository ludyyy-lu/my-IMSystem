package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"my-IMSystem/ws-gateway/internal/logic"
	"my-IMSystem/ws-gateway/internal/svc"
	"my-IMSystem/ws-gateway/internal/types"
)

func connectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConnectReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewConnectLogic(r.Context(), svcCtx)
		resp, err := l.Connect(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
