package svc

import (
	"log"
	"my-IMSystem/common/kafka"
	"my-IMSystem/friend-service/internal/config"
	"my-IMSystem/friend-service/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	// 仅迁移 friend-service 自己维护的表，避免改写 users 表结构。
	err = db.AutoMigrate(&model.FriendRequest{}, &model.Friend{}, &model.BlockedUser{})
	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}
	kafka.InitKafkaProducer(c.Kafka.Brokers)
	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
