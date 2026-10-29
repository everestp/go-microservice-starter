package logx

import (
	"log/slog"
	"os"
	"strings"
)

// New creates a JSON logger configured for production style output.
func New(service, env, level string) *slog.Logger {
	lvl := slog.LevelInfo
	switch strings.ToUpper(level) {
	case "DEBUG":
		lvl = slog.LevelDebug
	case "WARN", "WARNING":
		lvl = slog.LevelWarn
	case "ERROR":
		lvl = slog.LevelError
	}

	opts := &slog.HandlerOptions{Level: lvl}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler).With("service", service, "env", env)
	return logger
}
