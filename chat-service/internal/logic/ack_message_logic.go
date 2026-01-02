package logic

import (
	"context"

	chat_chat "my-IMSystem/chat-service/chat"
	"my-IMSystem/chat-service/internal/model"
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

// AckMessage 处理消息回执逻辑
// 1. 接收消息 ID 和用户 ID
// 2. 调用 model 层函数更新消息状态为已读
// 3. 返回处理结果
// 注意：这里的逻辑是同步的，可能会影响性能，实际应用中可以考虑异步处理或使用消息队列
func (l *AckMessageLogic) AckMessage(in *chat_chat.AckMessageReq) (*chat_chat.AckMessageResp, error) {
	// todo: add your logic here and delete this line
	err := model.MarkMessageAsRead(l.svcCtx.DB, in.MessageId, in.UserId)
	if err != nil {
		l.Logger.Errorf("Ack failed: %v", err)
		return &chat_chat.AckMessageResp{
			Status: "FAILED",
		}, nil
	}

	return &chat_chat.AckMessageResp{
		Status: "OK",
	}, nil
	// return &chat_chat.AckMessageResp{}, nil
}
