package rpc

import (
	"context"
	"errors"

	"my-IMSystem/auth-service/auth"
)

type AuthService struct {
	client auth.AuthClient
}

func NewAuthService(client auth.AuthClient) *AuthService {
	return &AuthService{client: client}
}

func (s *AuthService) VerifyToken(ctx context.Context, token string) (int64, error) {
	if s == nil || s.client == nil {
		return 0, errors.New("auth client not initialized")
	}
	resp, err := s.client.VerifyToken(ctx, &auth.VerifyTokenReq{AccessToken: token})
	if err != nil {
		return 0, err
	}
	if !resp.Valid {
		return 0, errors.New("invalid token")
	}
	return resp.UserId, nil
}
