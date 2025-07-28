package svc

import (
	"log"
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

	// 自动建表
	err = db.AutoMigrate(&model.FriendRequest{}, &model.Friend{})
	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
	}
}
