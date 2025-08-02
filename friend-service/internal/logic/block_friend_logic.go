package logic

import (
	"context"
	"errors"
	"time"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/model"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type BlockFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBlockFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BlockFriendLogic {
	return &BlockFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BlockFriendLogic) BlockFriend(in *friend_friend.BlockFriendReq) (*friend_friend.BlockFriendResp, error) {
	// todo: add your logic here and delete this line
	userID := in.UserId
	targetID := in.TargetId

	// 1. 防止重复拉黑（幂等）
	var existing model.BlockedUser
	err := l.svcCtx.DB.Where("user_id = ? AND blocked_id = ?", userID, targetID).First(&existing).Error
	if err == nil {
		// 已经拉黑过了
		return &friend_friend.BlockFriendResp{Msg: "已拉黑该用户"}, nil
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, status.Error(codes.Internal, "查询拉黑记录失败")
	}

	// 2. 插入拉黑记录
	block := model.BlockedUser{
		UserID:    userID,
		BlockedID: targetID,
		CreatedAt: time.Now(),
	}

	if err := l.svcCtx.DB.Create(&block).Error; err != nil {
		return nil, status.Error(codes.Internal, "插入拉黑记录失败")
	}

	return &friend_friend.BlockFriendResp{Msg: "拉黑成功"}, nil
}
