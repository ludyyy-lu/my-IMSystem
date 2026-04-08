package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
zrpc.RpcServerConf
Mysql struct {
DataSource string
}
Kafka struct {
Brokers []string
Topic   string
Group   string
// ReadReceiptTopic 是已读回执事件的 Kafka topic。
// chat-service 在标记消息已读后向此 topic 发布事件，
// ws-gateway 订阅后将回执实时推送给消息发送方。
ReadReceiptTopic string
}
FriendRpc zrpc.RpcClientConf
}
