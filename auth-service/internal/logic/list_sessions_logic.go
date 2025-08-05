package logic

import (
	"context"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/model"
	"my-IMSystem/auth-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ListSessionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSessionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSessionsLogic {
	return &ListSessionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 多设备：列出当前用户的所有活跃会话
func (l *ListSessionsLogic) ListSessions(in *auth_auth.ListSessionsReq) (*auth_auth.ListSessionsResp, error) {
	var sessions []model.Session
	err := l.svcCtx.DB.Where("uid = ?", in.UserId).Find(&sessions).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "查询会话失败: "+err.Error())
	}

	var respSessions []*auth_auth.Session
	for _, s := range sessions {
		respSessions = append(respSessions, &auth_auth.Session{
			DeviceId:  s.DeviceId,
			LoginAt:   s.CreatedAt.Unix(),
			ExpiresAt: s.ExpiresAt.Unix(),
		})
	}

	return &auth_auth.ListSessionsResp{
		Sessions: respSessions,
	}, nil
}
