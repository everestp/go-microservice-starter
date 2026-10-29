package service

import (
	"context"
	"sort"
	"time"

	"github.com/example/go-microservice-starter/pkg/clock"
	"github.com/example/go-microservice-starter/services/notes-service/internal/domain"
	"github.com/google/uuid"
)

// Repository persists notes.
type Repository interface {
	Create(ctx context.Context, note domain.Note) error
	Get(ctx context.Context, id string) (domain.Note, error)
	List(ctx context.Context) ([]domain.Note, error)
	Update(ctx context.Context, note domain.Note) error
	Delete(ctx context.Context, id string) error
}

// NotesService implements note CRUD and optimistic concurrency.
type NotesService struct {
	repo  Repository
	clock clock.Clock
}

func NewNotesService(repo Repository, c clock.Clock) *NotesService {
	return &NotesService{repo: repo, clock: c}
}

func (s *NotesService) Create(ctx context.Context, title, body string) (domain.Note, error) {
	now := s.clock.Now()
	note := domain.Note{
		ID:        uuid.NewString(),
		Title:     title,
		Body:      body,
		Version:   1,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, note); err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (s *NotesService) Get(ctx context.Context, id string) (domain.Note, error) {
	return s.repo.Get(ctx, id)
}

func (s *NotesService) List(ctx context.Context, cursor string, limit int) ([]domain.Note, string, error) {
	notes, err := s.repo.List(ctx)
	if err != nil {
		return nil, "", err
	}
	sort.Slice(notes, func(i, j int) bool {
		if notes[i].UpdatedAt.Equal(notes[j].UpdatedAt) {
			return notes[i].ID < notes[j].ID
		}
		return notes[i].UpdatedAt.After(notes[j].UpdatedAt)
	})
	if limit <= 0 || limit > len(notes) {
		limit = len(notes)
	}
	if cursor != "" {
		for i, item := range notes {
			if item.ID == cursor {
				if i+1 < len(notes) {
					notes = notes[i+1:]
				} else {
					notes = nil
				}
				break
			}
		}
	}
	if limit < len(notes) {
		out := notes[:limit]
		if len(notes) > limit {
			return out, notes[limit-1].ID, nil
		}
		return out, "", nil
	}
	return notes, "", nil
}

func (s *NotesService) Update(ctx context.Context, id string, version int64, title, body string) (domain.Note, error) {
	note, err := s.repo.Get(ctx, id)
	if err != nil {
		return domain.Note{}, err
	}
	if note.Version != version {
		return domain.Note{}, domain.ErrNoteVersionMismatch
	}
	note.Title = title
	note.Body = body
	note.Version++
	note.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, note); err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (s *NotesService) Delete(ctx context.Context, id string) error {
	_, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
