package svc

import (
	"log"
	"my-IMSystem/user-service/internal/config"
	"my-IMSystem/user-service/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 DB
	dsn := c.Mysql.DataSource // 从配置文件读取
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.User{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
	}

	// model.InitDB(c.Mysql.DataSource) //解耦调用
	// return &ServiceContext{
	// 	Config: c,
	// }
}
