package logic

import (
	"context"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendRequestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendRequestsLogic {
	return &GetFriendRequestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFriendRequestsLogic) GetFriendRequests(in *friend_friend.GetFriendRequestsRequest) (*friend_friend.GetFriendRequestsResponse, error) {
	// todo: add your logic here and delete this line

	return &friend_friend.GetFriendRequestsResponse{}, nil
}
