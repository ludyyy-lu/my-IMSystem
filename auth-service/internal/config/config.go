package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	JwtAuth struct {
		AccessSecret []byte
		AccessExpire int64
	}
	Mysql struct {
		DataSource string
	}
}
