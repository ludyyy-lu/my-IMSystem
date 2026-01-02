package logic

import (
	"context"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/model"
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
	// 1. 查找好友列表
	var friendRelations []model.Friend
	if err := l.svcCtx.DB.Where("user_id = ?", in.UserId).Find(&friendRelations).Error; err != nil {
		return nil, err
	}
	// 2. 提取好友ID
	friendIDs := make([]int64, 0)
	for _, fr := range friendRelations {
		friendIDs = append(friendIDs, fr.FriendID)
	}
	// 3. 查询好友用户名（跨表，需要访问用户表）
	var users []model.User // 假设你有个 User 模型
	if err := l.svcCtx.DB.Where("id IN ?", friendIDs).Find(&users).Error; err != nil {
		return nil, err
	}
	// 4. 封装结果
	friendInfos := make([]*friend_friend.FriendInfo, 0)
	for _, u := range users {
		friendInfos = append(friendInfos, &friend_friend.FriendInfo{
			FriendId: u.ID,
			Username: u.Username,
		})
	}

	return &friend_friend.GetFriendsResponse{
		Friends: friendInfos,
	}, nil
}
