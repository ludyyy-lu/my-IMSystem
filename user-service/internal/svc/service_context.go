package svc

import (
	"log"
	"my-IMSystem/user-service/internal/config"
	"my-IMSystem/user-service/internal/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB         // 数据库连接
	UserModel   *model.UserModel // 数据库访问接口
	RedisClient *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {

	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	// 自动建表
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatalf("auto migration failed: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB, // 选择默认数据库
	})

	return &ServiceContext{
		Config:      c,
		DB:          db,
		UserModel:   model.NewUserModel(db),
		RedisClient: rdb,
	}
}
