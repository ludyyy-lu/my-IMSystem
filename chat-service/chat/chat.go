package chat

import (
"context"

"github.com/zeromicro/go-zero/zrpc"
"google.golang.org/grpc"
)

// Chat is the go-zero client-side interface for the Chat gRPC service.
type (
Chat interface {
SendMessage(ctx context.Context, in *SendMessageReq, opts ...grpc.CallOption) (*SendMessageResp, error)
GetChatHistory(ctx context.Context, in *GetChatHistoryReq, opts ...grpc.CallOption) (*GetChatHistoryResp, error)
AckMessage(ctx context.Context, in *AckMessageReq, opts ...grpc.CallOption) (*AckMessageResp, error)
BatchAckMessages(ctx context.Context, in *BatchAckReq, opts ...grpc.CallOption) (*BatchAckResp, error)
GetUnreadCount(ctx context.Context, in *GetUnreadCountReq, opts ...grpc.CallOption) (*GetUnreadCountResp, error)
}

defaultChat struct {
cli zrpc.Client
}
)

func NewChat(cli zrpc.Client) Chat {
return &defaultChat{cli: cli}
}

func (m *defaultChat) SendMessage(ctx context.Context, in *SendMessageReq, opts ...grpc.CallOption) (*SendMessageResp, error) {
client := NewChatClient(m.cli.Conn())
return client.SendMessage(ctx, in, opts...)
}

func (m *defaultChat) GetChatHistory(ctx context.Context, in *GetChatHistoryReq, opts ...grpc.CallOption) (*GetChatHistoryResp, error) {
client := NewChatClient(m.cli.Conn())
return client.GetChatHistory(ctx, in, opts...)
}

func (m *defaultChat) AckMessage(ctx context.Context, in *AckMessageReq, opts ...grpc.CallOption) (*AckMessageResp, error) {
client := NewChatClient(m.cli.Conn())
return client.AckMessage(ctx, in, opts...)
}

func (m *defaultChat) BatchAckMessages(ctx context.Context, in *BatchAckReq, opts ...grpc.CallOption) (*BatchAckResp, error) {
client := NewChatClient(m.cli.Conn())
return client.BatchAckMessages(ctx, in, opts...)
}

func (m *defaultChat) GetUnreadCount(ctx context.Context, in *GetUnreadCountReq, opts ...grpc.CallOption) (*GetUnreadCountResp, error) {
client := NewChatClient(m.cli.Conn())
return client.GetUnreadCount(ctx, in, opts...)
}
