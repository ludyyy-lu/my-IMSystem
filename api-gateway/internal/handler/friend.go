package handler

import (
	"encoding/json"
	"net/http"

	"my-IMSystem/api-gateway/internal/middleware"
	"my-IMSystem/api-gateway/internal/svc"
	"my-IMSystem/api-gateway/internal/types"
	friendpb "my-IMSystem/friend-service/friend"
)

func decodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// GetFriendsHandler handles GET /api/friends
func GetFriendsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		resp, err := svcCtx.FriendRpc.GetFriends(r.Context(), &friendpb.GetFriendsRequest{UserId: userID})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "success", resp.Friends)
	}
}

// GetFriendRequestsHandler handles GET /api/friends/requests
func GetFriendRequestsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		resp, err := svcCtx.FriendRpc.GetFriendRequests(r.Context(), &friendpb.GetFriendRequestsRequest{UserId: userID})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "success", resp.Requests)
	}
}

// SendFriendRequestHandler handles POST /api/friends/request
func SendFriendRequestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		var req types.SendFriendRequestReq
		if err := decodeJSON(r, &req); err != nil {
			respondJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		resp, err := svcCtx.FriendRpc.SendFriendRequest(r.Context(), &friendpb.SendFriendRequestRequest{
			FromUserId: userID,
			ToUserId:   req.ToUserID,
			Remark:     req.Remark,
		})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, resp.Message, nil)
	}
}

// RespondFriendRequestHandler handles PUT /api/friends/request
func RespondFriendRequestHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RespondFriendRequestReq
		if err := decodeJSON(r, &req); err != nil {
			respondJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		accept := req.Action == "accept"
		resp, err := svcCtx.FriendRpc.RespondFriendRequest(r.Context(), &friendpb.RespondFriendRequestRequest{
			RequestId: req.RequestID,
			Accept:    accept,
		})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, resp.Message, nil)
	}
}

// DeleteFriendHandler handles DELETE /api/friends/:id
func DeleteFriendHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		friendIDStr := r.PathValue("id")
		friendID, err := parseInt64(friendIDStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, "invalid friend id", nil)
			return
		}
		resp, err := svcCtx.FriendRpc.DeleteFriend(r.Context(), &friendpb.DeleteFriendRequest{
			UserId:   userID,
			FriendId: friendID,
		})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, resp.Message, nil)
	}
}
