package logic

import (
	"context"
	"fmt"

	"my-IMSystem/pkg/jwt"
	"my-IMSystem/user-service/internal/svc"
	"my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/metadata"
)

type GetProfileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProfileLogic) GetProfile(in *user.GetProfileRequest) (*user.GetProfileResponse, error) {
	// Step 1: 从 metadata 中提取 token
	md, ok := metadata.FromIncomingContext(l.ctx)
	if !ok || len(md["authorization"]) == 0 {
		return nil, fmt.Errorf("缺少 token")
	}
	token := md["authorization"][0]

	// Step 2: 解析 token 拿到 uid
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return nil, fmt.Errorf("无效 token: %v", err)
	}
	uid := claims.Uid

	// Step 3: 查数据库
	var userModel struct {
		ID       int64  `gorm:"column:id"`
		Username string `gorm:"column:username"`
	}
	if err := l.svcCtx.DB.Table("users").Where("id = ?", uid).First(&userModel).Error; err != nil {
		return nil, err
	}

	// Step 4: 构造响应
	return &user.GetProfileResponse{
		Uid:      userModel.ID,
		Username: userModel.Username,
		Avatar:   "https://i.imgtg.com/2023/11/30/avatar.png", // 先写死
	}, nil
}
