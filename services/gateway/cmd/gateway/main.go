package main

import (
	"fmt"
	"log"
	"os"

	"github.com/example/go-microservice-starter/pkg/logx"
	"github.com/example/go-microservice-starter/services/gateway/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("gateway config: %v", err)
	}
	logger := logx.New(cfg.ServiceName, cfg.AppEnv, cfg.LogLevel)
	logger.Info("gateway service starting", "http_addr", cfg.HTTPAddr, "grpc_addr", cfg.GRPCAddr)
	fmt.Fprintf(os.Stdout, "service=%s http=%s grpc=%s\n", cfg.ServiceName, cfg.HTTPAddr, cfg.GRPCAddr)
}
