package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	JwtSecret    string
	AuthRpcConf  zrpc.RpcClientConf
	UserRpcConf  zrpc.RpcClientConf
	FriendRpcConf zrpc.RpcClientConf
	ChatRpcConf  zrpc.RpcClientConf
}
