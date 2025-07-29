package svc

import (
	"my-IMSystem/ws-gateway/internal/config"
	"my-IMSystem/ws-gateway/internal/conn"
)

type ServiceContext struct {
	Config      config.Config
	ConnManager *conn.ConnManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:      c,
		ConnManager: conn.NewConnManager(),
	}
}
