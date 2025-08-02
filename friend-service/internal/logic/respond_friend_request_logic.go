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
	// 1. 查询请求是否存在
	var fr model.FriendRequest
	if err := l.svcCtx.DB.First(&fr, in.RequestId).Error; err != nil {
		return nil, errors.New("好友请求不存在")
	}

	if fr.Status != "pending" {
		return nil, errors.New("该好友请求已处理过")
	}

	if in.Accept {
		// 2. 更新请求状态为 accepted
		fr.Status = "accepted"
		if err := l.svcCtx.DB.Save(&fr).Error; err != nil {
			return nil, err
		}

		// 3. 创建好友关系（双向插入）  // 原子性
		friendPairs := []model.Friend{
			{UserID: fr.FromUserID, FriendID: fr.ToUserID},
			{UserID: fr.ToUserID, FriendID: fr.FromUserID},
		}
		if err := l.svcCtx.DB.Create(&friendPairs).Error; err != nil {
			return nil, err
		}

		if err := kafka.SendMessage("im-friend-topic", common_model.FriendEvent{
			EventType: "request_accepted",
			FromUser:  fr.FromUserID,
			ToUser:    fr.ToUserID,
			Timestamp: time.Now().Unix(),
			Extra:     fr.Remark,
		}); err != nil {
			return nil, errors.New("failed to send friend request event to Kafka: " + err.Error())
		}

		return &friend_friend.RespondFriendRequestResponse{Message: "已接受好友请求"}, nil
	} else {
		// 拒绝：更新请求状态为 rejected
		fr.Status = "rejected"
		if err := l.svcCtx.DB.Save(&fr).Error; err != nil {
			return nil, err
		}
		if err := kafka.SendMessage("im-friend-topic", common_model.FriendEvent{
			EventType: "request_rejected",
			FromUser:  fr.FromUserID,
			ToUser:    fr.ToUserID,
			Timestamp: time.Now().Unix(),
			Extra:     fr.Remark,
		}); err != nil {
			return nil, errors.New("failed to send friend request event to Kafka: " + err.Error())
		}
		return &friend_friend.RespondFriendRequestResponse{Message: "已拒绝好友请求"}, nil
	}
}
