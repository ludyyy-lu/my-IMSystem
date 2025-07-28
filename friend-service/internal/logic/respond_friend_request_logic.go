package logic

import (
	"context"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RespondFriendRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRespondFriendRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RespondFriendRequestLogic {
	return &RespondFriendRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RespondFriendRequestLogic) RespondFriendRequest(in *friend_friend.RespondFriendRequestRequest) (*friend_friend.RespondFriendRequestResponse, error) {
	// todo: add your logic here and delete this line

	return &friend_friend.RespondFriendRequestResponse{}, nil
}
