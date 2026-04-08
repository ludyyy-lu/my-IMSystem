package svc

import (
	"my-IMSystem/api-gateway/internal/config"
	"my-IMSystem/auth-service/auth"
	"my-IMSystem/chat-service/chat"
	"my-IMSystem/friend-service/friend"
	"my-IMSystem/user-service/user"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	AuthRpc   auth.AuthClient
	UserRpc   user.UserClient
	FriendRpc friend.FriendClient
	ChatRpc   chat.ChatClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		AuthRpc:   auth.NewAuthClient(zrpc.MustNewClient(c.AuthRpcConf).Conn()),
		UserRpc:   user.NewUserClient(zrpc.MustNewClient(c.UserRpcConf).Conn()),
		FriendRpc: friend.NewFriendClient(zrpc.MustNewClient(c.FriendRpcConf).Conn()),
		ChatRpc:   chat.NewChatClient(zrpc.MustNewClient(c.ChatRpcConf).Conn()),
	}
}
