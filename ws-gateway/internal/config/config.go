package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Redis struct {
        Addr     string
        Password string
        DB       int
    }

    Kafka struct {
        Brokers []string
        Topic   string
    }
}
