package svc

import (
	"my-IMSystem/user-service/internal/config"
	"my-IMSystem/user-service/internal/model"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	model.InitDB()
	return &ServiceContext{
		Config: c,
	}
}
