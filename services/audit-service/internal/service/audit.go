package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/example/go-microservice-starter/pkg/clock"
	"github.com/example/go-microservice-starter/services/audit-service/internal/domain"
	"github.com/google/uuid"
)

// Repository persists append-only audit events.
type Repository interface {
	Save(ctx context.Context, event domain.AuditEvent) error
	List(ctx context.Context) ([]domain.AuditEvent, error)
}

// AuditService records and reads audit activity.
type AuditService struct {
	repo  Repository
	clock clock.Clock
}

func NewAuditService(repo Repository, c clock.Clock) *AuditService {
	return &AuditService{repo: repo, clock: c}
}

func (s *AuditService) Record(ctx context.Context, userID, eventType, source string, data map[string]string) (domain.AuditEvent, error) {
	event := domain.AuditEvent{
		ID:        uuid.NewString(),
		Type:      eventType,
		Version:   1,
		Source:    source,
		CreatedAt: s.clock.Now(),
		UserID:    userID,
		Data:      data,
	}
	if err := s.repo.Save(ctx, event); err != nil {
		return domain.AuditEvent{}, err
	}
	return event, nil
}

func (s *AuditService) List(ctx context.Context, cursor string, limit int) ([]domain.AuditEvent, error) {
	events, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit > len(events) {
		limit = len(events)
	}
	if cursor != "" {
		filtered := make([]domain.AuditEvent, 0, len(events))
		for _, event := range events {
			if event.ID == cursor {
				continue
			}
			filtered = append(filtered, event)
		}
		events = filtered
	}
	if limit < len(events) {
		events = events[:limit]
	}
	return events, nil
}

func (s *AuditService) MarshalData(data map[string]string) string {
	if len(data) == 0 {
		return "{}"
	}
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("{\"error\":%q}", err.Error())
	}
	return string(b)
}

func (s *AuditService) ValidatedAt(ts time.Time) bool {
	return !ts.IsZero()
}
