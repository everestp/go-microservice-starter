package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	AppEnv      string
	ServiceName string
	LogLevel    string
	HTTPAddr    string
	GRPCAddr    string
	DatabaseURL string
	NATSURL     string
	Shutdown    time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:      getenv("APP_ENV", "local"),
		ServiceName: getenv("SERVICE_NAME", "audit-service"),
		LogLevel:    getenv("LOG_LEVEL", "INFO"),
		HTTPAddr:    getenv("HTTP_ADDR", ":8083"),
		GRPCAddr:    getenv("GRPC_ADDR", ":9093"),
		DatabaseURL: getenv("DATABASE_URL", "postgres://starter:starter@localhost:5432/audit_db"),
		NATSURL:     getenv("NATS_URL", "nats://localhost:4222"),
	}
	shutdown, err := time.ParseDuration(getenv("SHUTDOWN_TIMEOUT", "15s"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg.Shutdown = shutdown
	return cfg, nil
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}
