package logic

import (
	"context"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/model"
	"my-IMSystem/auth-service/internal/svc"
	"my-IMSystem/common/utils"
	"my-IMSystem/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登录：验证账号密码，返回 token
func (l *LoginLogic) Login(in *auth_auth.LoginReq) (*auth_auth.LoginResp, error) {
	// todo: add your logic here and delete this line
	// 查用户
	var user model.User
	if err := l.svcCtx.DB.Where("username = ?", in.Username).First(&user).Error; err != nil {
		return nil, err
	}

	// 验证密码
	if !utils.CheckPasswordHash(in.Password, user.Password) {
		return nil, status.Error(codes.Unauthenticated, "用户名或密码错误")
	}

	// 生成 token
	accessToken, err := jwt.GenerateToken(user.ID, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, err
	}

	// 保存 Session
	session := &model.Session{
		Uid:          user.ID,
		DeviceId:     in.DeviceId,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		// ExpiresAt:    accessExpire,
	}
	if err := l.svcCtx.DB.Create(session).Error; err != nil {
		return nil, err
	}

	// 返回 token 信息
	return &auth_auth.LoginResp{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		// ExpiresAt:    accessExpire.Unix(),
	}, nil
}
