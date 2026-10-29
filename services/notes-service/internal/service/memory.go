package service

import (
	"context"
	"sync"

	"github.com/example/go-microservice-starter/services/notes-service/internal/domain"
)

// MemoryRepository is a simple in-memory implementation for local development and tests.
type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]domain.Note
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{items: map[string]domain.Note{}}
}

func (r *MemoryRepository) Create(_ context.Context, note domain.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[note.ID]; exists {
		return domain.ErrNoteNotFound
	}
	r.items[note.ID] = note
	return nil
}

func (r *MemoryRepository) Get(_ context.Context, id string) (domain.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	note, ok := r.items[id]
	if !ok {
		return domain.Note{}, domain.ErrNoteNotFound
	}
	return note, nil
}

func (r *MemoryRepository) List(_ context.Context) ([]domain.Note, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]domain.Note, 0, len(r.items))
	for _, item := range r.items {
		out = append(out, item)
	}
	return out, nil
}

func (r *MemoryRepository) Update(_ context.Context, note domain.Note) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[note.ID]; !exists {
		return domain.ErrNoteNotFound
	}
	r.items[note.ID] = note
	return nil
}

func (r *MemoryRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.items[id]; !exists {
		return domain.ErrNoteNotFound
	}
	delete(r.items, id)
	return nil
}
