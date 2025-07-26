package logic

import (
	"context"

	"my-IMSystem/user-service/internal/svc"
	"my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user_user.RegisterRequest) (*user_user.RegisterResponse, error) {
	// todo: add your logic here and delete this line

	return &user_user.RegisterResponse{}, nil
}
