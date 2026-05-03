package config

import (
	"fmt"
	"os"
)

type (
	Config struct {
		GRPC
		Postgres
	}

	GRPC struct {
		Port        string `env:"GRPC_PORT"`
		GatewayPort string `env:"GRPC_GATEWAY_PORT"`
	}

	Postgres struct {
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		DB       string `env:"POSTGRES_DB"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{
		GRPC: GRPC{
			Port:        getEnv("GRPC_PORT", "8080"),
			GatewayPort: getEnv("GRPC_GATEWAY_PORT", "8081"),
		},
		Postgres: Postgres{
			Host:     getEnv("POSTGRES_HOST", "localhost"),
			Port:     getEnv("POSTGRES_PORT", "5432"),
			User:     getEnv("POSTGRES_USER", "postgres"),
			Password: getEnv("POSTGRES_PASSWORD", "postgres"),
			DB:       getEnv("POSTGRES_DB", "library"),
		},
	}

	return cfg, nil
}

func (p Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DB)
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
