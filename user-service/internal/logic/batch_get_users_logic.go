package logic

import (
	"context"

	"my-IMSystem/user-service/internal/svc"
	user_user "my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	users, err := l.svcCtx.UserModel.FindByIDs(l.ctx, in.Uids)
	if err != nil {
		return nil, status.Error(codes.Internal, "获取用户信息失败")
	}

	var userInfos []*user_user.UserInfo
	for _, u := range users {
		userInfos = append(userInfos, &user_user.UserInfo{
			Id:        u.ID,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
			Bio:       u.Bio,
			CreatedAt: u.CreatedAt.Unix(),
			Disabled:  u.Disabled,
			Gender:    u.Gender,
		})
	}

	return &user_user.BatchGetUsersResp{
		Users: userInfos,
	}, nil
}
