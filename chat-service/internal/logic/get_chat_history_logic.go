package logic

import (
	"context"

	"my-IMSystem/chat-service/chat"
	"my-IMSystem/chat-service/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatHistoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatHistoryLogic {
	return &GetChatHistoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetChatHistoryLogic) GetChatHistory(in *chat_chat.GetChatHistoryReq) (*chat_chat.GetChatHistoryResp, error) {
	// todo: add your logic here and delete this line

	return &chat_chat.GetChatHistoryResp{}, nil
}
