package svc

import (
	"my-IMSystem/auth-service/auth"
	"my-IMSystem/chat-service/chat"
	"my-IMSystem/common/kafka"
	"my-IMSystem/ws-gateway/internal/config"
	"my-IMSystem/ws-gateway/internal/push"
	"my-IMSystem/ws-gateway/internal/rpc"
	"my-IMSystem/ws-gateway/internal/session"

	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config         config.Config
	SessionManager *session.Manager
	OfflineStore   session.OfflineStore
	PresenceStore  session.PresenceStore
	RedisClient    *redis.Client
	ChatRpc        chat.ChatClient
	AuthService    *rpc.AuthService
	PushService    *push.Service
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
	})

	kafka.InitKafkaProducer(c.Kafka.Brokers)

	sessionManager := session.NewManager()
	offlineStore := session.NewRedisOfflineMsgStore(rdb)
	presenceStore := session.NewRedisPresenceStore(rdb)
	chatClient := chat.NewChatClient(zrpc.MustNewClient(c.ChatRpcConf).Conn())
	authClient := auth.NewAuthClient(zrpc.MustNewClient(c.AuthRpcConf).Conn())

	return &ServiceContext{
		Config:         c,
		SessionManager: sessionManager,
		OfflineStore:   offlineStore,
		PresenceStore:  presenceStore,
		RedisClient:    rdb,
		ChatRpc:        chatClient,
		AuthService:    rpc.NewAuthService(authClient),
		PushService:    push.NewService(sessionManager, offlineStore),
	}
}
