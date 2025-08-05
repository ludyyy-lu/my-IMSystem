package logic

import (
	"context"

	"my-IMSystem/user-service/internal/svc"
	user_user "my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetUsersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetUsersLogic {
	return &BatchGetUsersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetUsersLogic) BatchGetUsers(in *user_user.BatchGetUsersReq) (*user_user.BatchGetUsersResp, error) {
	// todo: add your logic here and delete this line

	return &user_user.BatchGetUsersResp{}, nil
}
