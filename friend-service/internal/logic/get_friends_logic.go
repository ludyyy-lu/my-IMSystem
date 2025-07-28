package logic

import (
	"context"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendsLogic {
	return &GetFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendsLogic) GetFriends(in *friend_friend.GetFriendsRequest) (*friend_friend.GetFriendsResponse, error) {
	// todo: add your logic here and delete this line

	return &friend_friend.GetFriendsResponse{}, nil
}
