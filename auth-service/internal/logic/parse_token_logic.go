package logic

import (
	"context"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"
	"my-IMSystem/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
)

type ParseTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewParseTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ParseTokenLogic {
	return &ParseTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 解析 token，提取 userId（如 userId, deviceId）
func (l *ParseTokenLogic) ParseToken(in *auth_auth.ParseTokenReq) (*auth_auth.ParseTokenResp, error) {
	claims, err := jwt.ParseToken(in.AccessToken, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}

	return &auth_auth.ParseTokenResp{
		UserId:    claims.Uid,
		DeviceId:  claims.DeviceId,
		ExpiresAt: claims.ExpiresAt.Time.Unix(), // 从 jwt.RegisteredClaims 拿的
	}, nil
	// return &auth_auth.ParseTokenResp{}, nil
}
