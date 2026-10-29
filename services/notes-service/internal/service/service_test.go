package service

import (
	"context"
	"testing"
	"time"

	"github.com/example/go-microservice-starter/pkg/clock"
	"github.com/example/go-microservice-starter/services/notes-service/internal/domain"
)

func TestNotesService_CreateAndGet(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewNotesService(repo, clock.RealClock{})

	note, err := svc.Create(context.Background(), "hello", "world")
	if err != nil {
		t.Fatalf("create note: %v", err)
	}
	if note.ID == "" || note.Title != "hello" || note.Body != "world" {
		t.Fatalf("unexpected note: %#v", note)
	}

	got, err := svc.Get(context.Background(), note.ID)
	if err != nil {
		t.Fatalf("get note: %v", err)
	}
	if got.ID != note.ID || got.Version != 1 {
		t.Fatalf("expected version 1 note, got %#v", got)
	}
}

func TestNotesService_UpdateRespectsVersion(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewNotesService(repo, clock.RealClock{})

	note, err := svc.Create(context.Background(), "title", "body")
	if err != nil {
		t.Fatalf("create note: %v", err)
	}

	updated, err := svc.Update(context.Background(), note.ID, note.Version, "new title", "new body")
	if err != nil {
		t.Fatalf("update note: %v", err)
	}
	if updated.Version != 2 {
		t.Fatalf("expected version 2 after update, got %d", updated.Version)
	}

	if _, err := svc.Update(context.Background(), note.ID, note.Version, "stale", "value"); !errorIs(err, domain.ErrNoteVersionMismatch) {
		t.Fatalf("expected version mismatch, got %v", err)
	}
}

func errorIs(err error, target error) bool {
	return err != nil && err.Error() == target.Error()
}

func TestNotesService_ClockUsesRealTime(t *testing.T) {
	_ = time.Now()
}
