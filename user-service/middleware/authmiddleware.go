package middleware

import (
	"context"
	"net/http"
	"strings"

	"my-IMSystem/pkg/jwt" // 引入你刚写的 jwt 工具包
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 格式应该是 Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			logx.Errorf("Token parse error: %v", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// 将 uid 存入 context 中，供后续逻辑使用
		ctx := context.WithValue(r.Context(), "uid", claims.Uid)
		next(w, r.WithContext(ctx))
	}
}
