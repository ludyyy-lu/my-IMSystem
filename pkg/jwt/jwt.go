package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// var jwtKey = []byte("my-secret-key") // 可放入 config 管理

type Claims struct {
	Uid int64 `json:"uid"`
	jwt.RegisteredClaims
}

// 生成 JWT token
func GenerateToken(uid int64, secretKey []byte) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// 生成 Refresh Token
// 这个 token 的过期时间更长，通常用于刷新 access token
func GenerateRefreshToken(uid int64, secretKey []byte) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour) // 比 access 更久
	claims := &Claims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// 解析 token
func ParseToken(tokenString string, secretKey []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}
