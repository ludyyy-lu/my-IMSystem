package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	chat_chat "my-IMSystem/chat-service/chat"
	"my-IMSystem/chat-service/internal/svc"
	"my-IMSystem/friend-service/friend"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/status"
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
	// Step 0: 检查是否被对方拉黑
	blockResp, err := l.svcCtx.FriendRpc.IsBlocked(l.ctx, &friend.IsBlockedReq{
		SenderId: in.FromUserId, // 发消息的人
		TargetId: in.ToUserId,   // 收消息的人
	})
	if err != nil {
		return nil, fmt.Errorf("check block status failed: %w", err)
	}
	if blockResp.IsBlocked {
		return nil, status.Error(403, "你已被对方拉黑，无法发送消息")
	}
	// 1. 生成消息 ID + 时间戳
	msgID := uuid.New().String()
	timestamp := time.Now().Unix()

	// 2. 封装消息结构
	message := map[string]interface{}{
		"message_id": msgID,
		"from":       in.FromUserId,
		"to":         in.ToUserId,
		"content":    in.Content,
		"timestamp":  timestamp,
	}

	// 3. 编码为 JSON
	msgBytes, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("marshal message failed: %w", err)
	}
	// 4. 发送到 Kafka
	err = l.svcCtx.Producer.SendMessage(strconv.FormatInt(in.ToUserId, 10), msgBytes)
	// err = l.svcCtx.Producer.SendMessage(in.ToUserId.String(), msgBytes)
	if err != nil {
		return nil, fmt.Errorf("send kafka failed: %w", err)
	}

	return &chat_chat.SendMessageResp{
		Status:    "OK",
		MessageId: msgID,
		Timestamp: timestamp,
	}, nil
}
