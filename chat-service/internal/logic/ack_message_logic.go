package logic

import (
	"context"

	"my-IMSystem/chat-service/chat"
	"my-IMSystem/chat-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type AckMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAckMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AckMessageLogic {
	return &AckMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AckMessageLogic) AckMessage(in *chat_chat.AckMessageReq) (*chat_chat.AckMessageResp, error) {
	// todo: add your logic here and delete this line

	return &chat_chat.AckMessageResp{}, nil
}
