package logic

import (
	"context"

	"my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 生成 token（用于其他服务调用）
func (l *GenerateTokenLogic) GenerateToken(in *auth_auth.GenerateTokenReq) (*auth_auth.GenerateTokenResp, error) {
	// todo: add your logic here and delete this line

	return &auth_auth.GenerateTokenResp{}, nil
}
