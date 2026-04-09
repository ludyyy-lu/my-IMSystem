package handler

import (
	"net/http"
	"path"
	"strconv"

	"my-IMSystem/api-gateway/internal/middleware"
	"my-IMSystem/api-gateway/internal/svc"
	chatpb "my-IMSystem/chat-service/chat"
)

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// GetChatHistoryHandler handles GET /api/messages/:peer_id
func GetChatHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		peerIDStr := path.Base(r.URL.Path)
		peerID, err := parseInt64(peerIDStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, "invalid peer_id", nil)
			return
		}
		limitStr := r.URL.Query().Get("limit")
		limit := int64(50)
		if limitStr != "" {
			if n, e := parseInt64(limitStr); e == nil {
				limit = n
			}
		}
		resp, err := svcCtx.ChatRpc.GetChatHistory(r.Context(), &chatpb.GetChatHistoryReq{
			UserId: userID,
			PeerId: peerID,
			Limit:  limit,
		})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "success", resp.Messages)
	}
}

// GetUnreadCountHandler handles GET /api/messages/unread
func GetUnreadCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		resp, err := svcCtx.ChatRpc.GetUnreadCount(r.Context(), &chatpb.GetUnreadCountReq{UserId: userID})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "success", resp.Conversations)
	}
}
