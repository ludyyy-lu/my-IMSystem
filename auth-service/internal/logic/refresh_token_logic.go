package logic

import (
	"context"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"

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
	// todo: add your logic here and delete this line

	return &auth_auth.RefreshTokenResp{}, nil
}
