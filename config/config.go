package config

import "os"

type (
	Config struct {
		GRPC
	}

	GRPC struct {
		Port        string `env:"GRPC_PORT"`
		GatewayPort string `env:"GRPC_GATEWAY_PORT"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{
		GRPC{
			Port:        getEnv("GRPC_PORT", "8080"),
			GatewayPort: getEnv("GRPC_GATEWAY_PORT", "8081"),
		},
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
