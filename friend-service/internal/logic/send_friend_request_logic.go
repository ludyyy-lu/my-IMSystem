package logic

import (
	"context"

	"my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendFriendRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendFriendRequestLogic {
	return &SendFriendRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendFriendRequestLogic) SendFriendRequest(in *friend_friend.SendFriendRequestRequest) (*friend_friend.SendFriendRequestResponse, error) {
	// todo: add your logic here and delete this line

	return &friend_friend.SendFriendRequestResponse{}, nil
}
