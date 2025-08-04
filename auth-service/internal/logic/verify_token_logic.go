package logic

import (
	"context"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyTokenLogic {
	return &VerifyTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 校验 token 是否有效
func (l *VerifyTokenLogic) VerifyToken(in *auth_auth.VerifyTokenReq) (*auth_auth.VerifyTokenResp, error) {
	// todo: add your logic here and delete this line

	return &auth_auth.VerifyTokenResp{}, nil
}
