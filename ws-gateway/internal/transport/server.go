package transport

import (
	"net/http"

	"my-IMSystem/ws-gateway/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func Register(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/ws/connect",
				Handler: ConnectHandler(serverCtx),
			},
		},
	)
}
