package main

import (
	"flag"
	"fmt"

	"my-IMSystem/ws-gateway/internal/config"
	"my-IMSystem/ws-gateway/internal/consume"
	"my-IMSystem/ws-gateway/internal/svc"
	"my-IMSystem/ws-gateway/internal/transport"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/ws.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	// 启动好友事件消费者
	stopConsumers := consume.StartConsumersWithCancel(c.Kafka.Brokers, c.Kafka.Topic, "im-friend-topic", ctx.PushService)
	defer stopConsumers()
	transport.Register(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
