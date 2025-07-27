package logic

import (
	"context"

	"my-IMSystem/user-service/internal/model"
	"my-IMSystem/user-service/internal/svc"
	"my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
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

func (l *GetProfileLogic) GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	// 1. 从 ctx 中取出 uid（我们等会会把它注入进来）
	uidVal := ctx.Value("uid")
	if uidVal == nil {
		return nil, status.Error(401, "未授权访问")
	}
	uid := uidVal.(int64)

	// 2. 查询数据库
	var u model.User
	if err := l.svcCtx.DB.Where("id = ?", uid).First(&u).Error; err != nil {
		return nil, status.Error(500, "查询用户失败")
	}

	// 3. 返回响应
	return &user.GetProfileResponse{
		Uid:      u.ID,
		Username: u.Username,
	}, nil
}

// func (l *GetProfileLogic) GetProfile(in *user.GetProfileRequest) (*user.GetProfileResponse, error) {
// 	// todo: add your logic here and delete this line
// 	return &user.GetProfileResponse{}, nil
// }
