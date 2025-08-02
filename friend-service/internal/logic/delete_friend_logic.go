package logic

import (
	"context"
	"errors"
	"time"

	"my-IMSystem/common/common_model"
	"my-IMSystem/common/kafka"
	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/model"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DeleteFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFriendLogic {
	return &DeleteFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFriendLogic) DeleteFriend(in *friend_friend.DeleteFriendRequest) (*friend_friend.DeleteFriendResponse, error) {
	// todo: add your logic here and delete this line
	// 双向删除好友关系
	// gorm.DB.Transaction() 确保原子性
	err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND friend_id = ?", in.UserId, in.FriendId).
			Delete(&model.Friend{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ? AND friend_id = ?", in.FriendId, in.UserId).
			Delete(&model.Friend{}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 发送 Kafka 消息
	if err := kafka.SendMessage("im-friend-topic", common_model.FriendEvent{
		EventType: "friend_deleted",
		FromUser:  in.UserId,
		ToUser:    in.FriendId,
		Timestamp: time.Now().Unix(),
	}); err != nil {
		return nil, errors.New("failed to send friend deletion event to Kafka: " + err.Error())
	}
	return &friend_friend.DeleteFriendResponse{
		Message: "删除好友成功",
	}, nil
}
