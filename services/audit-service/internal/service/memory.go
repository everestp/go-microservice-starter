package service

import (
	"context"
	"sync"

	"github.com/example/go-microservice-starter/services/audit-service/internal/domain"
)

// MemoryRepository stores audit events in memory for local development and tests.
type MemoryRepository struct {
	mu    sync.RWMutex
	items []domain.AuditEvent
	byID  map[string]domain.AuditEvent
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{items: []domain.AuditEvent{}, byID: map[string]domain.AuditEvent{}}
}

func (r *MemoryRepository) Save(_ context.Context, event domain.AuditEvent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byID[event.ID]; exists {
		return nil
	}
	r.byID[event.ID] = event
	r.items = append(r.items, event)
	return nil
}

func (r *MemoryRepository) List(_ context.Context) ([]domain.AuditEvent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := append([]domain.AuditEvent(nil), r.items...)
	return out, nil
}
