// Package transport handles the HTTP→WebSocket upgrade and authentication.
// Once a connection is established it:
//   - verifies the bearer token via AuthService,
//   - marks the user online in the PresenceStore,
//   - registers the session with the SessionManager,
//   - starts session I/O goroutines,
//   - asynchronously delivers any pending offline messages.
//
// On connection close the user is marked offline and the session is removed
// from the SessionManager.  Inbound message routing is delegated to the router
// package; this file contains no business logic.
package transport

import (
	"context"
	"net/http"
	"strings"

	"my-IMSystem/ws-gateway/internal/router"
	"my-IMSystem/ws-gateway/internal/session"
	"my-IMSystem/ws-gateway/internal/svc"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins; restrict in production if needed
	},
}

// ConnectHandler returns an http.HandlerFunc that upgrades the request to a
// WebSocket connection after verifying the bearer token.
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

		userID, err := svcCtx.AuthService.VerifyToken(r.Context(), token)
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
			logx.Errorf("WebSocket upgrade failed for user %d: %v", userID, err)
			return
		}

		// Mark the user online in the shared presence store so that other
		// gateway nodes and downstream services can observe their status.
		if svcCtx.PresenceStore != nil {
			if pErr := svcCtx.PresenceStore.SetOnline(context.Background(), userID); pErr != nil {
				logx.Errorf("failed to set online presence for user %d: %v", userID, pErr)
			}
		}

		onMessage := func(uid int64, payload []byte) {
			router.HandleMessage(svcCtx, uid, payload)
		}
		onClose := func(uid int64) {
			svcCtx.SessionManager.Remove(uid)
			// Mark the user offline when their WebSocket connection is closed.
			if svcCtx.PresenceStore != nil {
				if pErr := svcCtx.PresenceStore.SetOffline(context.Background(), uid); pErr != nil {
					logx.Errorf("failed to clear online presence for user %d: %v", uid, pErr)
				}
			}
			logx.Infof("WebSocket connection closed for user %d", uid)
		}

		sess := session.NewSession(userID, conn, onMessage, onClose)
		svcCtx.SessionManager.Add(sess)
		sess.Start()
		go deliverOfflineMessages(userID, svcCtx, sess)

		logx.Infof("WebSocket connection established for user %d", userID)
	}
}

// deliverOfflineMessages loads any pending messages from the offline store and
// pushes them directly to the newly connected session.  Each stored entry is
// already a serialised PushMessage envelope, so no re-serialisation is needed.
func deliverOfflineMessages(userID int64, svcCtx *svc.ServiceContext, sess *session.Session) {
	if svcCtx.OfflineStore == nil {
		return
	}
	messages, err := svcCtx.OfflineStore.LoadAndDelete(userID)
	if err != nil {
		logx.Errorf("failed to load offline messages for user %d: %v", userID, err)
		return
	}
	for _, data := range messages {
		if err := sess.Send(data); err != nil {
			logx.Errorf("failed to deliver offline message to user %d: %v", userID, err)
		}
	}
	if len(messages) > 0 {
		logx.Infof("delivered %d offline message(s) to user %d", len(messages), userID)
	}
}

// extractToken reads the bearer token from Authorization header, query string,
// or Sec-WebSocket-Protocol header (in that priority order).
// It also returns the raw Sec-WebSocket-Protocol value so it can be echoed back
// to satisfy browser WebSocket clients that use the protocol field for auth.
func extractToken(r *http.Request) (token, protocolToken string) {
	token = r.Header.Get("Authorization")
	if token == "" {
		token = r.URL.Query().Get("token")
	}
	protocolToken = r.Header.Get("Sec-WebSocket-Protocol")
	if token == "" {
		token = protocolToken
	}
	return strings.TrimPrefix(token, "Bearer "), protocolToken
}
