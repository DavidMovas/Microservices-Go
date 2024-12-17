package internal

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     string `env:"PORT" envDefault:"8030"`
	RPCPort  string `env:"RPC_PORT" envDefault:"5030"`
	GRPCPort string `env:"GRPC_PORT" envDefault:"5031"`
	EsURL    string `env:"ELASTICSEARCH_URL"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
