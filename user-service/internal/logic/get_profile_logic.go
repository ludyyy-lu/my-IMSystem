package logic

import (
	"context"

	"my-IMSystem/user-service/internal/svc"
	user_user "my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProfileLogic) GetProfile(in *user_user.GetProfileReq) (*user_user.GetProfileResp, error) {
	// todo: add your logic here and delete this line

	return &user_user.GetProfileResp{}, nil
}
