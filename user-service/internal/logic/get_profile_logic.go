package logic

import (
	"context"

	"my-IMSystem/user-service/internal/svc"
	user_user "my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
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

func (l *GetProfileLogic) GetProfile(in *user_user.GetProfileReq) (*user_user.GetProfileResp, error) {
	// todo: add your logic here and delete this line
	u, err := l.svcCtx.UserModel.FindByID(l.ctx, in.Uid)
	if err != nil {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	return &user_user.GetProfileResp{
		User: &user_user.UserInfo{
			Id:        u.ID,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
			Bio:       u.Bio,
			CreatedAt: u.CreatedAt.Unix(),
			Disabled:  u.Disabled,
			Gender:    u.Gender,
		},
	}, nil
}
