package svc

import (
	"my-IMSystem/user-service/internal/config"
	"my-IMSystem/user-service/internal/model"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	model.InitDB(c.Mysql.DataSource) //解耦调用
	return &ServiceContext{
		Config: c,
	}
}
