package logic

import (
	"context"
	"errors"

	friend_friend "my-IMSystem/friend-service/friend"
	"my-IMSystem/friend-service/internal/model"
	"my-IMSystem/friend-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type IsBlockedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsBlockedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsBlockedLogic {
	return &IsBlockedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsBlockedLogic) IsBlocked(in *friend_friend.IsBlockedReq) (*friend_friend.IsBlockedResp, error) {
	// todo: add your logic here and delete this line
	var bu model.BlockedUser
	err := l.svcCtx.DB.Where("user_id = ? AND blocked_id = ?", in.TargetId, in.SenderId).First(&bu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &friend_friend.IsBlockedResp{IsBlocked: false}, nil
	}
	if err != nil {
		l.Logger.Errorf("查询拉黑失败: %v", err)
		return nil, status.Error(500, "服务异常")
	}
	return &friend_friend.IsBlockedResp{IsBlocked: true}, nil
}
