package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	//  HTTP REST 服务 的配置结构 包含了 HTTP 服务监听地址、端口、超时、限流等参数
	rest.RestConf
	Redis struct {
		Addr     string
		Password string
		DB       int
	}

	Kafka struct {
		Brokers []string
		Topic   string
	}
	//  gRPC 客户端 的配置结构 包含服务端点（Endpoints）、超时时间、连接池等参数
	ChatRpcConf zrpc.RpcClientConf
	AuthRpcConf zrpc.RpcClientConf
}
