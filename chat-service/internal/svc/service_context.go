package svc

import (
	"log"
	"my-IMSystem/chat-service/internal/config"
	"my-IMSystem/chat-service/internal/handler"
	"my-IMSystem/chat-service/internal/kafka"
	"my-IMSystem/chat-service/internal/model"
	"my-IMSystem/friend-service/friend"

	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//ServiceContext 是 go-zero 框架用来做依赖注入的地方，它的作用是：
// 把全局需要共享的资源（DB、Redis、Kafka 等）集中在一起初始化，然后传递给每一个业务逻辑 handler 使用

type ServiceContext struct {
	Config       config.Config
	Producer     *kafka.KafkaProducer
	DB           *gorm.DB
	MessageModel model.MessageModel
	FriendRpc    friend.Friend
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化数据库
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	// 自动迁移（自动建表）
	if err := db.AutoMigrate(&model.Message{}); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}
	// 启动 Kafka 消费者
	kafka.StartChatMessageConsumer(
		c.Kafka.Brokers,
		c.Kafka.Topic,
		c.Kafka.Group,
		handler.ChatMessageHandler(db),
	)

	return &ServiceContext{
		Config: c,
		Producer: kafka.NewKafkaProducer(
			c.Kafka.Brokers,
			c.Kafka.Topic,
		),
		DB:           db,
		MessageModel: model.NewMessageModel(db),
		FriendRpc:    friend.NewFriend(zrpc.MustNewClient(c.FriendRpc)),
	}
}
