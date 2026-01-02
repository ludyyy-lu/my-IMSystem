package logic

import (
	"context"
	"time"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"
	"my-IMSystem/pkg/jwt"

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
	tokenStr, err := jwt.GenerateToken(in.Uid, in.DeviceId, l.svcCtx.JwtSecretKey)
	if err != nil {
		return nil, err
	}
	expiredAt := time.Now().Add(7 * 24 * time.Hour).Unix() // 和 jwt 里设置的过期时间保持一致
	return &auth_auth.GenerateTokenResp{
		Token:     tokenStr,
		ExpiredAt: expiredAt,
	}, nil
}
