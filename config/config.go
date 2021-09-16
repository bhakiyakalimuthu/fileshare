package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

type Config struct {
	SourceFilePath string `env:"FILE_PATH" envDefault:"."`
	DestFilePath   string `env:"FILE_PATH" envDefault:"./dest"`
	Host           string `env:"HOST" envDefault:"localhost"`
	GrpcServerPort int    `env:"GRPC_SERVER_PORT" envDefault:"1997"`
	HttpServerPort string `env:"HTTP_SERVER_PORT" envDefault:"8080"`
}

var cfg = &Config{}

func InitConfig() {
	if err := env.Parse(cfg); err != nil {
		panic(fmt.Sprintf("failed to initialize the default config %s", err))
	}
}

func Get() *Config {
	return cfg
}
