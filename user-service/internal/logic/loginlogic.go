package logic

import (
	"context"

	"my-IMSystem/user-service/internal/model"
	"my-IMSystem/user-service/internal/svc"
	"my-IMSystem/user-service/user"
	"my-IMSystem/pkg/jwt"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
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

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	// todo: add your logic here and delete this line
	var u model.User
	// 1. 查找用户
	err := l.svcCtx.DB.Where("username = ?", in.Username).First(&u).Error
	if err != nil {
		return &user.LoginResponse{
			Message: "用户不存在",
		}, nil // 返回 nil 是为了保持 gRPC 通信成功，但 message 提示错误
	}
	// 2. 校验密码（bcrypt）
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)); err != nil {
		return &user.LoginResponse{
			Message: "密码错误",
		}, nil
	}

	token, err := jwt.GenerateToken(u.ID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "token 生成失败: %v", err)
	}

	return &user.LoginResponse{
		Token:   token,
		Message: "登录成功",
	}, nil
}
