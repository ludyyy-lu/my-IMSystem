package svc

import (
	"my-IMSystem/chat-service/internal/config"
	"my-IMSystem/chat-service/internal/kafka"
)

type ServiceContext struct {
	Config   config.Config
	Producer *kafka.KafkaProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Producer: kafka.NewKafkaProducer(
			[]string{"kafka:9092"},
			"im-chat-topic",
		),
	}
}
