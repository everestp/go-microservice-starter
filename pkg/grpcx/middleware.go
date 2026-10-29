package grpcx

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type loggerKey struct{}

// UnaryServerInterceptor extracts request metadata and adds it to the logger context.
func UnaryServerInterceptor(logger *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			logger = logger.With("request_id", md.Get("request-id"), "user_id", md.Get("user-id"))
		}
		ctx = context.WithValue(ctx, loggerKey{}, logger)
		return handler(ctx, req)
	}
}
