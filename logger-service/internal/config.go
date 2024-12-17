package internal

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port          string     `env:"PORT" envDefault:"8030"`
	RPCPort       string     `env:"RPC_PORT" envDefault:"5030"`
	GRPCPort      string     `env:"GRPC_PORT" envDefault:"5031"`
	MongoURL      string     `env:"MONGO_URL" envDefault:"mongodb://mongo:27017"`
	MongoUser     string     `env:"MONGO_USER"`
	MongoPassword string     `env:"MONGO_PASSWORD"`
	LokiConfig    LokiConfig `envPrefix:"LOKI_"`
}

type LokiConfig struct {
	Host     string `env:"HOST" envDefault:"http://loki:3100/loki/api/v1/push"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Labels   map[string]string
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, err
	}

	c.LokiConfig.Labels = make(map[string]string)
	c.LokiConfig.Labels["service"] = "logger-service"

	return &c, nil
}
