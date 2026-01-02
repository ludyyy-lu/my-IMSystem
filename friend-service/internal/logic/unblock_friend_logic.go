package logic

import (
	"context"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/model"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
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
	db := l.svcCtx.DB
	err := db.Where("user_id = ? AND blocked_id = ?", in.UserId, in.TargetId).Delete(&model.BlockedUser{}).Error
	if err != nil {
		l.Logger.Errorf("Failed to unblock user: %v", err)
		return nil, status.Error(500, "取消拉黑失败")
	}
	return &friend_friend.UnblockFriendResp{Msg: "取消拉黑成功"}, nil
}
