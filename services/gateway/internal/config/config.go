package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	AppEnv        string
	ServiceName   string
	LogLevel      string
	HTTPAddr      string
	GRPCAddr      string
	AuthGRPCAddr  string
	NotesGRPCAddr string
	AuditGRPCAddr string
	Shutdown      time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		AppEnv:        getenv("APP_ENV", "local"),
		ServiceName:   getenv("SERVICE_NAME", "gateway"),
		LogLevel:      getenv("LOG_LEVEL", "INFO"),
		HTTPAddr:      getenv("HTTP_ADDR", ":8080"),
		GRPCAddr:      getenv("GRPC_ADDR", ":9090"),
		AuthGRPCAddr:  getenv("AUTH_GRPC_ADDR", "auth-service:9090"),
		NotesGRPCAddr: getenv("NOTES_GRPC_ADDR", "notes-service:9090"),
		AuditGRPCAddr: getenv("AUDIT_GRPC_ADDR", "audit-service:9090"),
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
