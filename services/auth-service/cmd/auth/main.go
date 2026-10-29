package main

import (
	"log"
	"net"

	authv1 "github.com/example/go-microservice-starter/gen/proto/auth/v1"
	"github.com/example/go-microservice-starter/pkg/clock"
	"github.com/example/go-microservice-starter/pkg/logx"
	"github.com/example/go-microservice-starter/services/auth-service/internal/config"
	"github.com/example/go-microservice-starter/services/auth-service/internal/repository/postgres"
	"github.com/example/go-microservice-starter/services/auth-service/internal/service"
	grpcserver "github.com/example/go-microservice-starter/services/auth-service/internal/transport/grpc"
	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("auth config: %v", err)
	}
	logger := logx.New(cfg.ServiceName, cfg.AppEnv, cfg.LogLevel)
	repo := postgres.NewRepository()
	hasher := service.NewArgon2Hasher(65536, 3, 2)
	issuer := service.NewStaticTokenIssuer()
	svc := service.NewAuthService(repo, hasher, issuer, clock.RealClock{})
	grpcServer := grpcserver.NewServer(svc, logger)
	listener, err := net.Listen("tcp", cfg.GRPCAddr)
	if err != nil {
		log.Fatalf("listen grpc: %v", err)
	}
	server := grpc.NewServer()
	authv1.RegisterAuthServiceServer(server, grpcServer)
	logger.Info("auth service starting", "grpc_addr", cfg.GRPCAddr)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("serve grpc: %v", err)
	}
}
