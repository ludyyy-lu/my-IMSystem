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

type LogoutSessionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutSessionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutSessionLogic {
	return &LogoutSessionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 多设备：注销某个设备的 token（退出登录）
func (l *LogoutSessionLogic) LogoutSession(in *auth_auth.LogoutSessionReq) (*auth_auth.LogoutSessionResp, error) {
	// 尝试删除该用户某个设备的 Session
	err := l.svcCtx.DB.Where("uid = ? AND device_id = ?", in.UserId, in.DeviceId).Delete(&model.Session{}).Error
	if err != nil {
		return nil, status.Error(codes.Internal, "注销设备失败: "+err.Error())
	}

	return &auth_auth.LogoutSessionResp{
		Message: "设备登出成功",
	}, nil
}
