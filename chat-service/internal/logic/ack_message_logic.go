package logic

import (
"context"
"time"

chat_chat "my-IMSystem/chat-service/chat"
"my-IMSystem/chat-service/internal/model"
"my-IMSystem/chat-service/internal/svc"
"my-IMSystem/common/common_model"

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

// AckMessage 将单条消息标记为已读，并向 Kafka 发布已读回执事件通知消息发送方。
//
// 流程：
//  1. 调用 model.MarkMessageAsReadAndGet 更新 DB 消息状态（幂等）。
//  2. 若配置了 ReadReceiptTopic，向 Kafka 发布 ReadReceiptEvent。
func (l *AckMessageLogic) AckMessage(in *chat_chat.AckMessageReq) (*chat_chat.AckMessageResp, error) {
if in.MessageId == "" {
return &chat_chat.AckMessageResp{Status: "FAILED"}, nil
}

msg, err := model.MarkMessageAsReadAndGet(l.svcCtx.DB, in.MessageId, in.UserId)
if err != nil {
l.Logger.Errorf("AckMessage: mark failed msgID=%s userID=%d err=%v", in.MessageId, in.UserId, err)
return &chat_chat.AckMessageResp{Status: "FAILED"}, nil
}

// 发布已读回执事件（异步通知发送方，失败不影响主流程）
if l.svcCtx.ReadReceiptProducer != nil && msg != nil {
event := common_model.ReadReceiptEvent{
MessageId: in.MessageId,
SenderId:  msg.FromUserId,
ReaderId:  in.UserId,
ReadAt:    time.Now().UnixMilli(),
}
if pubErr := l.svcCtx.ReadReceiptProducer.SendMessage(in.MessageId, event); pubErr != nil {
l.Logger.Errorf("AckMessage: publish read receipt failed msgID=%s err=%v", in.MessageId, pubErr)
}
}

return &chat_chat.AckMessageResp{Status: "OK"}, nil
}
