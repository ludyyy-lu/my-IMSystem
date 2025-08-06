package config

import "github.com/zeromicro/go-zero/zrpc"
// 配置struct映射
type Config struct {
	zrpc.RpcServerConf
	Mysql struct {
		DataSource string
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
}
