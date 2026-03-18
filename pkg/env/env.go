package env

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

type Config struct {
	ServiceName string `env:"SERVICE_NAME" envDefault:"fleet"`
	HTTPPort    string `env:"HTTP_PORT" envDefault:"8081"`
	GRPCPort    string `env:"GRPC_PORT" envDefault:"9090"`
	DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://postgres:postgres@localhost:5432/school_tracking?sslmode=disable"`
	NatsURL     string `env:"NATS_URL" envDefault:"nats://localhost:4222"`
	Environment string `env:"ENVIRONMENT" envDefault:"development"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"debug"`
}

func NewConfig() *Config {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to load environment variables: %v", err)
	}
	return &cfg
}

var Module = fx.Module("env", fx.Provide(NewConfig))
