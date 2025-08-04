package logic

import (
	"context"

	"my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"

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
	// todo: add your logic here and delete this line

	return &auth_auth.ParseTokenResp{}, nil
}
