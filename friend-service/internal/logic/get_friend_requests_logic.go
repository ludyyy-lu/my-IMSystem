package logic

import (
	"context"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/model"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// 查看好友申请
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
	var requests []model.FriendRequest
	// 数据库查询
	err := l.svcCtx.DB.
		Where("to_user_id = ?", in.UserId).
		Order("created_at DESC").
		Find(&requests).Error
	if err != nil {
		return nil, err
	}
	// 构建响应
	var resp []*friend_friend.FriendRequest
	for _, r := range requests {
		resp = append(resp, &friend_friend.FriendRequest{
			RequestId:    r.ID,
			FromUserId:   r.FromUserID,
			FromUsername: "", // 暂时先不查用户名，后面可以关联 user-service 获取用户名
			Remark:       r.Remark,
			Status:       r.Status,
		})
	}

	return &friend_friend.GetFriendRequestsResponse{
		Requests: resp,
	}, nil
}
