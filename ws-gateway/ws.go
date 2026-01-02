package main

import (
	"flag"
	"fmt"

	"my-IMSystem/ws-gateway/internal/config"
	"my-IMSystem/ws-gateway/internal/conn"
	"my-IMSystem/ws-gateway/internal/handler"
	"my-IMSystem/ws-gateway/internal/kafka"
	"my-IMSystem/ws-gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ws-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()
	// 初始化连接池
	conn.InitGlobalConnManager()

	ctx := svc.NewServiceContext(c)
	// 启动好友事件消费者
	go kafka.StartFriendConsumer(c.Kafka.Brokers, "im-friend-topic")
	go kafka.StartChatConsumer(c.Kafka.Brokers, "im-chat-topic")
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
