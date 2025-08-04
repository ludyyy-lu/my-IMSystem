package logic

import (
	"context"

	"my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
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

	return &friend_friend.GetBlockedListResp{}, nil
}
