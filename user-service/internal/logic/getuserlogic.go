package logic

import (
	"context"
	"errors"

	"my-IMSystem/user-service/internal/model"
	"my-IMSystem/user-service/internal/svc"
	"my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserLogic {
	return &GetUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserLogic) GetUser(in *user.GetUserRequest) (*user.GetUserResponse, error) {
	// todo: add your logic here and delete this line
	// 根据uid查询用户信息
	var u model.User
	if err := l.svcCtx.DB.First(&u, in.Uid).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &user.GetUserResponse{
		Uid:      int64(u.ID),
		Username: u.Username,
	}, nil
}
