package logic

import (
	"context"
	"time"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"
	"my-IMSystem/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshTokenLogic {
	return &RefreshTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 刷新 token（用于保持会话）
func (l *RefreshTokenLogic) RefreshToken(in *auth_auth.RefreshTokenReq) (*auth_auth.RefreshTokenResp, error) {
	// 1. 解析 refresh_token
	claims, err := jwt.ParseToken(in.RefreshToken, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err // 非法或过期 token
	}

	// 2. 校验是否为 refresh_token（比如在 claims 里区分用途，或 Redis 判断是否还有效，以后再做）

	// 3. 生成新的 access_token（仍然用 device_id）
	accessToken, err := jwt.GenerateToken(claims.Uid, claims.DeviceId, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}

	// 4. 返回新 access_token + 过期时间戳
	expiresAt := time.Now().Add(7 * 24 * time.Hour).Unix() // 与 token 设置的过期时间保持一致
	return &auth_auth.RefreshTokenResp{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt,
	}, nil
}
