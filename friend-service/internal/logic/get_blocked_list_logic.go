package logic

import (
	"context"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/model"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
)

type GetBlockedListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetBlockedListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBlockedListLogic {
	return &GetBlockedListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetBlockedListLogic) GetBlockedList(in *friend_friend.GetBlockedListReq) (*friend_friend.GetBlockedListResp, error) {
	// todo: add your logic here and delete this line
	var list []model.BlockedUser
	err := l.svcCtx.DB.Where("user_id = ?", in.UserId).Find(&list).Error
	if err != nil {
		l.Logger.Errorf("查询拉黑列表失败: %v", err)
		return nil, status.Error(500, "获取拉黑列表失败")
	}
	// 提取 blocked_id 字段
	var ids []int64
	for _, bu := range list {
		ids = append(ids, bu.BlockedID)
	}

	return &friend_friend.GetBlockedListResp{
		BlockedIds: ids,
	}, nil
}
