package svc

import (
	"my-IMSystem/ws-gateway/internal/config"
	"my-IMSystem/ws-gateway/internal/conn"
)

type ServiceContext struct {
	Config       config.Config
	ConnManager  *conn.ConnManager
	OfflineStore *conn.OfflineMsgManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		ConnManager:  conn.NewConnManager(),
		OfflineStore: conn.NewOfflineMsgManager(),
	}
}
