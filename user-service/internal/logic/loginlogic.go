package logic

import (
	"context"
	"time"

	"my-IMSystem/user-service/internal/model"
	"my-IMSystem/user-service/internal/svc"
	"my-IMSystem/user-service/user"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
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

	// 3. 生成 JWT Token
	token, err := generateJWT(u.ID)
	if err != nil {
		return &user.LoginResponse{
			Message: "生成Token失败",
		}, nil
	}

	return &user.LoginResponse{
		Token:   token,
		Message: "登录成功",
	}, nil
}

func generateJWT(uid int64) (string, error) {
	claims := jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(24 * time.Hour).Unix(), // 有效期一天
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("secret123"))
}
