package logic

import (
	"context"
	"errors"

	auth_auth "my-IMSystem/auth-service/auth"
	"my-IMSystem/auth-service/internal/model"
	"my-IMSystem/auth-service/internal/svc"
	"my-IMSystem/pkg/jwt"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
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

// 用户注册：注册账号密码，返回 token
func (l *RegisterLogic) Register(in *auth_auth.RegisterReq) (*auth_auth.RegisterResp, error) {
	// todo: add your logic here and delete this line
	db := l.svcCtx.DB // GORM 实例

	// 1. 判断用户名是否已存在
	var existing model.User
	if err := db.Where("username = ?", in.Username).First(&existing).Error; err != gorm.ErrRecordNotFound {
		return nil, errors.New("用户名已存在")
	}

	// 2. 加密密码
	hashedPwd, err := utils.HashPassword(in.Password)
	if err != nil {
		return nil, err
	}

	// 3. 插入用户
	user := model.User{
		Username: in.Username,
		Password: hashedPwd,
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	// 4. 生成 token
	token, err := jwt.GenerateToken(user.ID, []byte(l.svcCtx.Config.JwtAuth.AccessSecret))
	if err != nil {
		return nil, err
	}

	return &auth_auth.RegisterResp{
		Token: token,
		Uid:   user.ID,
	}, nil
	// return &auth_auth.RegisterResp{}, nil
}
