package logic

import (
	"context"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/model"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

// 查看我发出的好友申请
type GetSentFriendRequestsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSentFriendRequestsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSentFriendRequestsLogic {
	return &GetSentFriendRequestsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSentFriendRequestsLogic) GetSentFriendRequests(in *friend_friend.GetFriendRequestsRequest) (*friend_friend.GetFriendRequestsResponse, error) {
	var requests []model.FriendRequest
	err := l.svcCtx.DB.
		Where("from_user_id = ?", in.UserId).
		Order("created_at DESC").
		Find(&requests).Error
	if err != nil {
		return nil, err
	}

	var resp []*friend_friend.FriendRequest
	for _, r := range requests {
		// The FriendRequest proto message is reused for sent requests.
		// FromUserId carries the recipient (ToUserID) so the caller knows who the request was sent to.
		// FromUsername (empty here) can be enriched by the caller if needed.
		resp = append(resp, &friend_friend.FriendRequest{
			RequestId:    r.ID,
			FromUserId:   r.ToUserID,
			FromUsername: "",
			Remark:       r.Remark,
			Status:       r.Status,
		})
	}

	return &friend_friend.GetFriendRequestsResponse{
		Requests: resp,
	}, nil
}
