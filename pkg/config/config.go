package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Config is the shared environment configuration loaded by each service.
type Config struct {
	AppEnv      string
	ServiceName string
	LogLevel    string
	GRPCAddr    string
	HTTPAddr    string
	DatabaseURL string
	NATSURL     string
	OTELURL     string
	Shutdown    time.Duration
}

// Load reads essential settings from environment variables and validates them.
func Load() (Config, error) {
	cfg := Config{
		AppEnv:      getEnv("APP_ENV", "local"),
		ServiceName: getEnv("SERVICE_NAME", "service"),
		LogLevel:    getEnv("LOG_LEVEL", "INFO"),
		GRPCAddr:    getEnv("GRPC_ADDR", ":9090"),
		HTTPAddr:    getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/app"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4222"),
		OTELURL:     getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4317"),
	}

	shutdown, err := time.ParseDuration(getEnv("SHUTDOWN_TIMEOUT", "15s"))
	if err != nil {
		return Config{}, fmt.Errorf("invalid SHUTDOWN_TIMEOUT: %w", err)
	}
	cfg.Shutdown = shutdown

	if cfg.ServiceName == "" {
		return Config{}, fmt.Errorf("SERVICE_NAME must not be empty")
	}
	if cfg.HTTPAddr == "" || cfg.GRPCAddr == "" {
		return Config{}, fmt.Errorf("HTTP_ADDR and GRPC_ADDR must be set")
	}
	if cfg.DatabaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL must not be empty")
	}
	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		return value
	}
	return fallback
}

func ParseInt(key string, fallback int) int {
	value := getEnv(key, strconv.Itoa(fallback))
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func ParseBool(key string, fallback bool) bool {
	value := getEnv(key, strconv.FormatBool(fallback))
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func SplitCSV(value string) []string {
	if value == "" {
		return nil
	}
	return strings.FieldsFunc(value, func(r rune) bool { return r == ',' || r == ';' })
}
