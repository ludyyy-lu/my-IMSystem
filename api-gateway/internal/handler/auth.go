package handler

import (
	"encoding/json"
	"net/http"

	"my-IMSystem/api-gateway/internal/svc"
	"my-IMSystem/api-gateway/internal/types"
	authpb "my-IMSystem/auth-service/auth"
)

func respondJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(types.Response{Code: code, Message: message, Data: data})
}

// RegisterHandler handles POST /api/auth/register
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		resp, err := svcCtx.AuthRpc.Register(r.Context(), &authpb.RegisterReq{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "注册成功", types.TokenResp{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    resp.ExpiresAt,
		})
	}
}

// LoginHandler handles POST /api/auth/login
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondJSON(w, http.StatusBadRequest, "invalid request body", nil)
			return
		}
		resp, err := svcCtx.AuthRpc.Login(r.Context(), &authpb.LoginReq{
			Username: req.Username,
			Password: req.Password,
			DeviceId: req.DeviceID,
		})
		if err != nil {
			respondJSON(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}
		respondJSON(w, http.StatusOK, "登录成功", types.TokenResp{
			AccessToken:  resp.AccessToken,
			RefreshToken: resp.RefreshToken,
			ExpiresAt:    resp.ExpiresAt,
		})
	}
}
