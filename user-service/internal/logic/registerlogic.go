package logic

import (
	"context"
	"errors"

	"my-IMSystem/user-service/internal/model"
	"my-IMSystem/user-service/internal/svc"
	"my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// todo: add your logic here and delete this line
	// Step 1: 检查用户名是否已存在
	var existing model.User
	if err := l.svcCtx.DB.Where("username = ?", in.Username).First(&existing).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}
	// Step 2: 加密密码
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}
	// Step 3: 存入数据库
	newUser := model.User{
		Username: in.Username,
		Password: string(hashedPwd),
	}

	if err := l.svcCtx.DB.Create(&newUser).Error; err != nil {
		return nil, errors.New("创建用户失败")
	}
	// Step 4: 返回响应
	return &user.RegisterResponse{
		Uid:     int64(newUser.ID),
		Message: "注册成功",
	}, nil
}
