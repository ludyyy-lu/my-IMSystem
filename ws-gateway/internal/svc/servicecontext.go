package svc

import (
	"my-IMSystem/ws-gateway/internal/config"
	"my-IMSystem/ws-gateway/internal/conn"

	"github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config       config.Config
	ConnManager  *conn.ConnManager
	OfflineStore *conn.RedisOfflineMsgStore // 离线消息存储
	RedisClient  *redis.Client // 如果需要 Redis 支持，可以添加
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // 有密码就填
        DB:       0,
    })
	
	return &ServiceContext{
		Config:       c,
		ConnManager:  conn.NewConnManager(),
		OfflineStore: conn.NewRedisOfflineMsgStore(rdb),
		RedisClient:  rdb,
	}
}
