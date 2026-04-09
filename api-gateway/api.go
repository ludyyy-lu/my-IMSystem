package main

import (
	"flag"
	"fmt"
	"net/http"

	"my-IMSystem/api-gateway/internal/config"
	"my-IMSystem/api-gateway/internal/handler"
	"my-IMSystem/api-gateway/internal/middleware"
	"my-IMSystem/api-gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/api.yaml", "the config file")

// corsMiddleware adds CORS headers to every response so the React dev server
// (or any browser client) can reach the gateway.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	secret := c.JwtSecret

	// Public routes (no auth required)
	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodPost,
			Path:    "/api/auth/register",
			Handler: handler.RegisterHandler(ctx),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/auth/login",
			Handler: handler.LoginHandler(ctx),
		},
	})

	// Authenticated routes
	auth := func(h http.HandlerFunc) http.HandlerFunc {
		return middleware.Auth(secret, h)
	}
	server.AddRoutes([]rest.Route{
		// User
		{
			Method:  http.MethodGet,
			Path:    "/api/users/me",
			Handler: auth(handler.GetMyProfileHandler(ctx)),
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/users/search",
			Handler: auth(handler.SearchUserHandler(ctx)),
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/users/:id",
			Handler: auth(handler.GetProfileHandler(ctx)),
		},
		{
			Method:  http.MethodPut,
			Path:    "/api/users/profile",
			Handler: auth(handler.UpdateProfileHandler(ctx)),
		},
		// Friends
		{
			Method:  http.MethodGet,
			Path:    "/api/friends",
			Handler: auth(handler.GetFriendsHandler(ctx)),
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/friends/requests",
			Handler: auth(handler.GetFriendRequestsHandler(ctx)),
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/friends/requests/sent",
			Handler: auth(handler.GetSentFriendRequestsHandler(ctx)),
		},
		{
			Method:  http.MethodPost,
			Path:    "/api/friends/request",
			Handler: auth(handler.SendFriendRequestHandler(ctx)),
		},
		{
			Method:  http.MethodPut,
			Path:    "/api/friends/request",
			Handler: auth(handler.RespondFriendRequestHandler(ctx)),
		},
		{
			Method:  http.MethodDelete,
			Path:    "/api/friends/:id",
			Handler: auth(handler.DeleteFriendHandler(ctx)),
		},
		// Messages
		{
			Method:  http.MethodGet,
			Path:    "/api/messages/unread",
			Handler: auth(handler.GetUnreadCountHandler(ctx)),
		},
		{
			Method:  http.MethodGet,
			Path:    "/api/messages/:peer_id",
			Handler: auth(handler.GetChatHistoryHandler(ctx)),
		},
	})

	// Wrap the entire server with CORS middleware
	server.Use(rest.ToMiddleware(corsMiddleware))

	fmt.Printf("Starting API gateway at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
