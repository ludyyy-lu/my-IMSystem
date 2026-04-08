package handler

import (
	"net/http"

	"my-IMSystem/api-gateway/internal/middleware"
	"my-IMSystem/api-gateway/internal/svc"
	"my-IMSystem/api-gateway/internal/types"
	userpb "my-IMSystem/user-service/user"
)

// GetProfileHandler handles GET /api/users/:id
// Path parameter :id is extracted from the URL path.
func GetProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		uid, err := parseInt64(idStr)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, "invalid user id", nil)
			return
		}
		resp, err := svcCtx.UserRpc.GetProfile(r.Context(), &userpb.GetProfileReq{Uid: uid})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "success", resp.User)
	}
}

// UpdateProfileHandler handles PUT /api/users/profile
func UpdateProfileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r)
		var req types.UpdateProfileReq
		if err := decodeJSON(r, &req); err != nil {
			respondJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		_, err := svcCtx.UserRpc.UpdateProfile(r.Context(), &userpb.UpdateProfileReq{
			UserId:   userID,
			Nickname: req.Nickname,
			Avatar:   req.Avatar,
			Bio:      req.Bio,
		})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "更新成功", nil)
	}
}

// SearchUserHandler handles GET /api/users/search?keyword=xxx
func SearchUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keyword := r.URL.Query().Get("keyword")
		resp, err := svcCtx.UserRpc.SearchUser(r.Context(), &userpb.SearchUserReq{Keyword: keyword})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "success", resp.Results)
	}
}
