package repository

import (
	"context"
	"sort"
	"sync"

	"github.com/example/go-microservice-starter/services/audit-service/internal/domain"
)

// MemoryRepository stores audit events in memory for local development and tests.
type MemoryRepository struct {
	mu    sync.RWMutex
	items []domain.AuditEvent
	byID  map[string]domain.AuditEvent
	index map[string]bool
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		items: []domain.AuditEvent{},
		byID:  map[string]domain.AuditEvent{},
		index: map[string]bool{},
	}
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
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.After(out[j].CreatedAt)
	})
	return out, nil
}
