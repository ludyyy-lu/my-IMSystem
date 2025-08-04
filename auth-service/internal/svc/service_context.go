package svc

import "my-IMSystem/auth-service/internal/config"

type ServiceContext struct {
	Config       config.Config
	JwtSecretKey []byte
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:       c,
		JwtSecretKey: []byte(c.JwtAuth.AccessSecret),
	}
}
