package svc

import (
	"my-IMSystem/chat-service/internal/config"
	"my-IMSystem/chat-service/internal/handler"
	"my-IMSystem/chat-service/internal/kafka"
)

//ServiceContext 是 go-zero 框架用来做依赖注入的地方，它的作用是：
// 把全局需要共享的资源（DB、Redis、Kafka 等）集中在一起初始化，然后传递给每一个业务逻辑 handler 使用

type ServiceContext struct {
	Config   config.Config
	Producer *kafka.KafkaProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
	kafka.StartChatMessageConsumer(
		[]string{"kafka:9092"},
		"im-chat-topic",
		"chat-group",
		handler.ChatMessageHandler,
	)

	return &ServiceContext{
		Config: c,
		Producer: kafka.NewKafkaProducer(
			[]string{"kafka:9092"},
			"im-chat-topic",
		),
	}
}
