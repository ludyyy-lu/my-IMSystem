package logic

import (
	"context"

	"my-IMSystem/chat-service/chat"
	"my-IMSystem/chat-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMessageLogic) SendMessage(in *chat_chat.SendMessageReq) (*chat_chat.SendMessageResp, error) {
	// todo: add your logic here and delete this line

	return &chat_chat.SendMessageResp{}, nil
}
