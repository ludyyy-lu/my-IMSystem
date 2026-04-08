package logic

import (
"context"
"fmt"
"time"

chat_chat "my-IMSystem/chat-service/chat"
"my-IMSystem/chat-service/internal/model"
"my-IMSystem/chat-service/internal/svc"
"my-IMSystem/common/common_model"

"github.com/zeromicro/go-zero/core/logx"
)

type BatchAckMessagesLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewBatchAckMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchAckMessagesLogic {
return &BatchAckMessagesLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// BatchAckMessages 将 peerId 发送给 userId 的所有未读消息一次性标记为已读，
// 并向 Kafka 发布批量已读回执事件。
//
// 适用场景：用户打开某个会话窗口时，一次性清零该会话的未读数。
// 幂等：重复调用无副作用。
func (l *BatchAckMessagesLogic) BatchAckMessages(in *chat_chat.BatchAckReq) (*chat_chat.BatchAckResp, error) {
if in.UserId == 0 || in.PeerId == 0 {
return &chat_chat.BatchAckResp{Status: "FAILED"}, nil
}

ackedCount, err := model.BatchMarkConversationAsRead(l.svcCtx.DB, in.UserId, in.PeerId)
if err != nil {
l.Logger.Errorf("BatchAckMessages: DB update failed userID=%d peerID=%d err=%v",
in.UserId, in.PeerId, err)
return &chat_chat.BatchAckResp{Status: "FAILED"}, nil
}

// 仅在有消息被更新时才发布回执事件
if ackedCount > 0 && l.svcCtx.ReadReceiptProducer != nil {
event := common_model.BatchReadReceiptEvent{
SenderId:   in.PeerId,
ReaderId:   in.UserId,
AckedCount: ackedCount,
ReadAt:     time.Now().UnixMilli(),
}
key := fmt.Sprintf("batch:%d:%d", in.PeerId, in.UserId)
if pubErr := l.svcCtx.ReadReceiptProducer.SendMessage(key, event); pubErr != nil {
l.Logger.Errorf("BatchAckMessages: publish batch receipt failed err=%v", pubErr)
}
}

return &chat_chat.BatchAckResp{Status: "OK", AckedCount: ackedCount}, nil
}
