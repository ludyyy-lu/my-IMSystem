package logic

import (
	"context"

	"my-IMSystem/user-service/internal/svc"
	user_user "my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsOnlineLogic {
	return &IsOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsOnlineLogic) IsOnline(in *user_user.IsOnlineReq) (*user_user.IsOnlineResp, error) {
	// todo: add your logic here and delete this line

	return &user_user.IsOnlineResp{}, nil
}
