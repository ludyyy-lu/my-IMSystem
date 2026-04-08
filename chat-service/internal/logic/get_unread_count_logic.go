package logic

import (
"context"

chat_chat "my-IMSystem/chat-service/chat"
"my-IMSystem/chat-service/internal/model"
"my-IMSystem/chat-service/internal/svc"

"github.com/zeromicro/go-zero/core/logx"
)

type GetUnreadCountLogic struct {
ctx    context.Context
svcCtx *svc.ServiceContext
logx.Logger
}

func NewGetUnreadCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnreadCountLogic {
return &GetUnreadCountLogic{
ctx:    ctx,
svcCtx: svcCtx,
Logger: logx.WithContext(ctx),
}
}

// GetUnreadCount 查询 userId 的未读消息统计，按会话对端聚合返回。
//
// 返回值：
//   - total_unread：所有会话未读消息总数
//   - conversations：每个有未读消息的会话详情（peer_id、未读数、最新消息预览）
func (l *GetUnreadCountLogic) GetUnreadCount(in *chat_chat.GetUnreadCountReq) (*chat_chat.GetUnreadCountResp, error) {
if in.UserId == 0 {
return &chat_chat.GetUnreadCountResp{}, nil
}

stats, err := model.GetUnreadCountByPeer(l.svcCtx.DB, in.UserId)
if err != nil {
l.Logger.Errorf("GetUnreadCount: query failed userID=%d err=%v", in.UserId, err)
return nil, err
}

var total int64
conversations := make([]*chat_chat.UnreadConversation, 0, len(stats))
for _, s := range stats {
total += s.UnreadCount
conversations = append(conversations, &chat_chat.UnreadConversation{
PeerId:        s.PeerId,
UnreadCount:   s.UnreadCount,
LastContent:   s.LastContent,
LastTimestamp: s.LastTimestamp.UnixMilli(),
})
}

return &chat_chat.GetUnreadCountResp{
TotalUnread:   total,
Conversations: conversations,
}, nil
}
