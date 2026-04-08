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

// ServiceContext 集中管理 chat-service 的全局依赖：数据库、Kafka 生产者、gRPC 客户端等。
type ServiceContext struct {
Config config.Config
// Producer 向聊天消息 topic 发布消息
Producer *kafka.KafkaProducer
// ReadReceiptProducer 向已读回执 topic 发布事件（nil 表示未配置，安全跳过）
ReadReceiptProducer *kafka.KafkaProducer
DB                  *gorm.DB
MessageModel        model.MessageModel
FriendRpc           friend.Friend
}

func NewServiceContext(c config.Config) *ServiceContext {
db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
if err != nil {
log.Fatalf("failed to connect DB: %v", err)
}
if err := db.AutoMigrate(&model.Message{}); err != nil {
log.Fatalf("failed to auto-migrate: %v", err)
}
kafka.StartChatMessageConsumer(
c.Kafka.Brokers,
c.Kafka.Topic,
c.Kafka.Group,
handler.ChatMessageHandler(db),
)

var readReceiptProducer *kafka.KafkaProducer
if c.Kafka.ReadReceiptTopic != "" {
readReceiptProducer = kafka.NewKafkaProducer(c.Kafka.Brokers, c.Kafka.ReadReceiptTopic)
}

return &ServiceContext{
Config:              c,
Producer:            kafka.NewKafkaProducer(c.Kafka.Brokers, c.Kafka.Topic),
ReadReceiptProducer: readReceiptProducer,
DB:                  db,
MessageModel:        model.NewMessageModel(db),
FriendRpc:           friend.NewFriend(zrpc.MustNewClient(c.FriendRpc)),
}
}
