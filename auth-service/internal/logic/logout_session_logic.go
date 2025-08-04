package logic

import (
	"context"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutSessionLogic {
	return &LogoutSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 多设备：注销某个设备的 token（退出登录）
func (l *LogoutSessionLogic) LogoutSession(in *auth_auth.LogoutSessionReq) (*auth_auth.LogoutSessionResp, error) {
	// todo: add your logic here and delete this line

	return &auth_auth.LogoutSessionResp{}, nil
}
