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
	db := l.svcCtx.DB
	// Step 0: 检查是否被拉黑
	// 添加好友可以直接查询好友库，但是chat服务就不行，微服务，所以用rpc跨服务调用
	// 此处直接查询数据库
	var blockedUser model.BlockedUser
	if err := db.Where("user_id = ? AND blocked_id = ?", in.ToUserId, in.FromUserId).First(&blockedUser).Error; err == nil {
		return nil, errors.New("对方已将你拉黑，无法发送好友请求")
	}

	// Step 1: 检查是否已经是好友
	var existingFriend model.Friend
	if err := db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		in.FromUserId, in.ToUserId, in.ToUserId, in.FromUserId).First(&existingFriend).Error; err == nil {
		return nil, errors.New("你们已经是好友了")
	}

	// Step 2: 检查是否已发送过请求
	var existingRequest model.FriendRequest
	if err := db.Where("from_user_id = ? AND to_user_id = ? AND status = ?", in.FromUserId, in.ToUserId, "pending").
		First(&existingRequest).Error; err == nil {
		return nil, errors.New("已经发送过好友请求了")
	}

	// Step 3: 创建好友请求
	request := model.FriendRequest{
		FromUserID: in.FromUserId,
		ToUserID:   in.ToUserId,
		Remark:     in.Remark,
		Status:     "pending",
	}
	if err := db.Create(&request).Error; err != nil {
		return nil, err
	}
	if err := kafka.SendMessage("im-friend-topic", common_model.FriendEvent{
		EventType: "request_sent",
		FromUser:  in.FromUserId,
		ToUser:    in.ToUserId,
		Timestamp: time.Now().Unix(),
		Extra:     in.Remark, // 附加信息
	}); err != nil {
		return nil, errors.New("failed to send friend request event to Kafka: " + err.Error())
	}
	return &friend_friend.SendFriendRequestResponse{
		Message: "好友请求已发送",
	}, nil
}
