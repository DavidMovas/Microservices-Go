package internal

import (
	"fmt"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8040"`
	Domain      string `env:"MAIL_DOMAIN"`
	Host        string `env:"MAIL_HOST"`
	MailPort    int    `env:"MAIL_PORT" envDefault:"1025"`
	Username    string `env:"MAIL_USERNAME"`
	Password    string `env:"MAIL_PASSWORD"`
	Encryption  string `env:"MAIL_ENCRYPTION"`
	FromAddress string `env:"MAIL_FROM_ADDRESS"`
	FromName    string `env:"MAIL_FROM_NAME"`
}

func NewConfig() (*Config, error) {
	_ = godotenv.Load()

	var c Config
	if err := env.Parse(&c); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &c, nil
}
