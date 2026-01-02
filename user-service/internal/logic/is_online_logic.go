package logic

import (
	"context"
	"fmt"

	"my-IMSystem/user-service/internal/svc"
	user_user "my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type IsOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewIsOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *IsOnlineLogic {
	return &IsOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *IsOnlineLogic) IsOnline(in *user_user.IsOnlineReq) (*user_user.IsOnlineResp, error) {
	key := "im:online:users"
	uid := in.Uid

	exists, err := l.svcCtx.RedisClient.SIsMember(l.ctx, key, uid).Result()
	if err != nil {
		return nil, fmt.Errorf("redis error: %v", err)
	}

	return &user_user.IsOnlineResp{
		Online: exists,
	}, nil
}
