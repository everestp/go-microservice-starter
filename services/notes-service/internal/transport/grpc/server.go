package grpc

import (
	"context"
	"log/slog"

	notesv1 "github.com/example/go-microservice-starter/gen/proto/notes/v1"
	"github.com/example/go-microservice-starter/services/notes-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server adapts the notes domain service to the generated protobuf contract.
type Server struct {
	notesv1.UnimplementedNotesServiceServer
	service *service.NotesService
	logger  *slog.Logger
}

func NewServer(svc *service.NotesService, logger *slog.Logger) *Server {
	return &Server{service: svc, logger: logger}
}

func (s *Server) Create(ctx context.Context, req *notesv1.CreateRequest) (*notesv1.CreateResponse, error) {
	note, err := s.service.Create(ctx, req.GetTitle(), req.GetBody())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &notesv1.CreateResponse{Id: note.ID, Title: note.Title, Body: note.Body, Version: note.Version}, nil
}

func (s *Server) List(ctx context.Context, req *notesv1.ListRequest) (*notesv1.ListResponse, error) {
	items, next, err := s.service.List(ctx, req.GetCursor(), int(req.GetLimit()))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &notesv1.ListResponse{NextCursor: next}
	for _, item := range items {
		resp.Items = append(resp.Items, &notesv1.Note{Id: item.ID, Title: item.Title, Body: item.Body, Version: item.Version})
	}
	return resp, nil
}

func (s *Server) Get(ctx context.Context, req *notesv1.GetRequest) (*notesv1.GetResponse, error) {
	note, err := s.service.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &notesv1.GetResponse{Note: &notesv1.Note{Id: note.ID, Title: note.Title, Body: note.Body, Version: note.Version}}, nil
}

func (s *Server) Update(ctx context.Context, req *notesv1.UpdateRequest) (*notesv1.UpdateResponse, error) {
	note, err := s.service.Update(ctx, req.GetId(), req.GetVersion(), req.GetTitle(), req.GetBody())
	if err != nil {
		return nil, mapError(err)
	}
	return &notesv1.UpdateResponse{Id: note.ID, Title: note.Title, Body: note.Body, Version: note.Version}, nil
}

func (s *Server) Delete(ctx context.Context, req *notesv1.DeleteRequest) (*notesv1.DeleteResponse, error) {
	if err := s.service.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &notesv1.DeleteResponse{}, nil
}

func mapError(err error) error {
	switch err {
	case nil:
		return nil
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
