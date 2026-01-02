package logic

import (
	"context"
	"strconv"
	"time"

	chat_chat "my-IMSystem/chat-service/chat"
	"my-IMSystem/chat-service/internal/model"
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

// 调用model层的函数处理请求
func (l *GetChatHistoryLogic) GetChatHistory(in *chat_chat.GetChatHistoryReq) (*chat_chat.GetChatHistoryResp, error) {
	// todo: add your logic here and delete this line
	// 1. 参数处理
	before := time.Now()
	if in.BeforeTimestamp > 0 {
		before = time.UnixMilli(in.BeforeTimestamp)
	}

	// 2. 查询数据库
	msgs, err := model.GetChatMessages(l.svcCtx.DB, in.UserId, in.PeerId, int(in.Limit), before)
	if err != nil {
		return nil, err
	}

	// 3. 转换成消息响应格式
	var resp []*chat_chat.ChatMessage
	for _, m := range msgs {
		resp = append(resp, &chat_chat.ChatMessage{
			FromUserId: m.FromUserId,
			ToUserId:   m.ToUserId,
			Content:    m.Content,
			Timestamp:  m.CreatedAt.UnixMilli(),
			MessageId:  strconv.FormatInt(m.Id, 10),
			IsRead:     m.Status == 1,
		})
	}
	return &chat_chat.GetChatHistoryResp{
		Messages: resp,
	}, nil
}
