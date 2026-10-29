package grpc

import (
	"context"
	"errors"
	"log/slog"

	authv1 "github.com/example/go-microservice-starter/gen/proto/auth/v1"
	"github.com/example/go-microservice-starter/services/auth-service/internal/domain"
	"github.com/example/go-microservice-starter/services/auth-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server adapts the auth domain service to the generated protobuf contract.
type Server struct {
	authv1.UnimplementedAuthServiceServer
	service *service.AuthService
	logger  *slog.Logger
}

// NewServer creates a grpc server adapter for the auth service.
func NewServer(svc *service.AuthService, logger *slog.Logger) *Server {
	return &Server{service: svc, logger: logger}
}

func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	pair, err := s.service.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, mapError(err)
	}
	return &authv1.RegisterResponse{UserId: pair.UserID, AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	pair, err := s.service.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, mapError(err)
	}
	return &authv1.LoginResponse{UserId: pair.UserID, AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Server) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	pair, err := s.service.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, mapError(err)
	}
	return &authv1.RefreshResponse{AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}

func (s *Server) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if err := s.service.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, mapError(err)
	}
	return &authv1.LogoutResponse{}, nil
}

func (s *Server) Me(ctx context.Context, _ *authv1.MeRequest) (*authv1.MeResponse, error) {
	userID, ok := ctx.Value("user_id").(string)
	if !ok || userID == "" {
		return nil, status.Error(codes.Unauthenticated, "missing user identity")
	}
	user, err := s.service.Me(ctx, userID)
	if err != nil {
		return nil, mapError(err)
	}
	return &authv1.MeResponse{UserId: user.ID, Email: user.Email}, nil
}

func mapError(err error) error {
	switch {
	case errors.Is(err, domain.ErrUserExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrUserNotFound), errors.Is(err, domain.ErrRefreshTokenNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrInvalidCredentials), errors.Is(err, domain.ErrInvalidRefreshToken), errors.Is(err, domain.ErrRefreshTokenRevoked):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
