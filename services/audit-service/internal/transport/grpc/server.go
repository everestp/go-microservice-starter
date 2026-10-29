package grpc

import (
	"context"
	"log/slog"

	auditv1 "github.com/example/go-microservice-starter/gen/proto/audit/v1"
	"github.com/example/go-microservice-starter/services/audit-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server adapts the audit domain service to the generated protobuf contract.
type Server struct {
	auditv1.UnimplementedAuditServiceServer
	service *service.AuditService
	logger  *slog.Logger
}

func NewServer(svc *service.AuditService, logger *slog.Logger) *Server {
	return &Server{service: svc, logger: logger}
}

func (s *Server) ListEvents(ctx context.Context, req *auditv1.ListEventsRequest) (*auditv1.ListEventsResponse, error) {
	events, err := s.service.List(ctx, req.GetCursor(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &auditv1.ListEventsResponse{}
	for _, event := range events {
		resp.Items = append(resp.Items, &auditv1.AuditEvent{
			Id:        event.ID,
			Type:      event.Type,
			Version:   event.Version,
			Source:    event.Source,
			CreatedAt: event.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UserId:    event.UserID,
			Data:      s.service.MarshalData(event.Data),
		})
	}
	return resp, nil
}
