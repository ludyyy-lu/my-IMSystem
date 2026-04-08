package server

import (
"context"

chat_chat "my-IMSystem/chat-service/chat"
"my-IMSystem/chat-service/internal/logic"
"my-IMSystem/chat-service/internal/svc"
)

type ChatServer struct {
svcCtx *svc.ServiceContext
chat_chat.UnimplementedChatServer
}

func NewChatServer(svcCtx *svc.ServiceContext) *ChatServer {
return &ChatServer{svcCtx: svcCtx}
}

func (s *ChatServer) SendMessage(ctx context.Context, in *chat_chat.SendMessageReq) (*chat_chat.SendMessageResp, error) {
l := logic.NewSendMessageLogic(ctx, s.svcCtx)
return l.SendMessage(in)
}

func (s *ChatServer) GetChatHistory(ctx context.Context, in *chat_chat.GetChatHistoryReq) (*chat_chat.GetChatHistoryResp, error) {
l := logic.NewGetChatHistoryLogic(ctx, s.svcCtx)
return l.GetChatHistory(in)
}

func (s *ChatServer) AckMessage(ctx context.Context, in *chat_chat.AckMessageReq) (*chat_chat.AckMessageResp, error) {
l := logic.NewAckMessageLogic(ctx, s.svcCtx)
return l.AckMessage(in)
}

func (s *ChatServer) BatchAckMessages(ctx context.Context, in *chat_chat.BatchAckReq) (*chat_chat.BatchAckResp, error) {
l := logic.NewBatchAckMessagesLogic(ctx, s.svcCtx)
return l.BatchAckMessages(in)
}

func (s *ChatServer) GetUnreadCount(ctx context.Context, in *chat_chat.GetUnreadCountReq) (*chat_chat.GetUnreadCountResp, error) {
l := logic.NewGetUnreadCountLogic(ctx, s.svcCtx)
return l.GetUnreadCount(in)
}
