package logic

import (
	"context"

	"my-IMSystem/user-service/internal/svc"
	user_user "my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SearchUserLogic) SearchUser(in *user_user.SearchUserReq) (*user_user.SearchUserResp, error) {
	users, err := l.svcCtx.UserModel.SearchByKeyword(l.ctx, in.Keyword)
	if err != nil {
		return nil, status.Error(codes.Internal, "搜索失败")
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
	return &user_user.SearchUserResp{
		Results: userInfos,
	}, nil
}
