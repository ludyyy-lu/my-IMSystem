package middleware

import (
	"context"
	"net/http"
	"strings"

	"my-IMSystem/pkg/jwt"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// Auth is an HTTP middleware that validates the JWT bearer token and injects
// the parsed user ID into the request context.
func Auth(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
		if token == "" {
			http.Error(w, `{"code":401,"message":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		claims, err := jwt.ParseToken(token, []byte(secret))
		if err != nil {
			http.Error(w, `{"code":401,"message":"invalid token"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, claims.Uid)
		next(w, r.WithContext(ctx))
	}
}

// GetUserID extracts the authenticated user ID from the request context.
func GetUserID(r *http.Request) int64 {
	if v := r.Context().Value(UserIDKey); v != nil {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}
