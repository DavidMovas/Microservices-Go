package internal

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port       string `env:"PORT" envDefault:"8020"`
	ConnString string `env:"POSTGRES_URL"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, err
	}

	return &c, nil
}
