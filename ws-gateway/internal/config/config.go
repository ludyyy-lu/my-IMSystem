package config

import (
"github.com/zeromicro/go-zero/rest"
"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
rest.RestConf
Redis struct {
Addr     string
Password string
DB       int
}
Kafka struct {
Brokers     []string
Topic       string
FriendTopic string
// ReadReceiptTopic 是已读回执事件的 Kafka topic
ReadReceiptTopic string
}
ChatRpcConf zrpc.RpcClientConf
AuthRpcConf zrpc.RpcClientConf
}
