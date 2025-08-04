package logic

import (
	"context"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSessionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSessionsLogic {
	return &ListSessionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 多设备：列出当前用户的所有活跃会话
func (l *ListSessionsLogic) ListSessions(in *auth_auth.ListSessionsReq) (*auth_auth.ListSessionsResp, error) {
	// todo: add your logic here and delete this line

	return &auth_auth.ListSessionsResp{}, nil
}
