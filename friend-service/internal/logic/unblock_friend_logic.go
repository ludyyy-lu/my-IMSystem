package logic

import (
	"context"

	"my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnblockFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnblockFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnblockFriendLogic {
	return &UnblockFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnblockFriendLogic) UnblockFriend(in *friend_friend.UnblockFriendReq) (*friend_friend.UnblockFriendResp, error) {
	// todo: add your logic here and delete this line

	return &friend_friend.UnblockFriendResp{}, nil
}
