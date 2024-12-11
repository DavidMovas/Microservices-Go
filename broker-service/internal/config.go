package internal

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string `env:"PORT" envDefault:"8010"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	fmt.Printf("PORT: %s\n", c.Port)

	return &c, nil
}
