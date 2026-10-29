package service

import (
	"context"
	"testing"
	"time"

	"github.com/example/go-microservice-starter/pkg/clock"
)

func TestAuditService_ListEvents(t *testing.T) {
	repo := NewMemoryRepository()
	svc := NewAuditService(repo, clock.RealClock{})

	if _, err := svc.Record(context.Background(), "user-1", "login", "auth", map[string]string{"status": "ok"}); err != nil {
		t.Fatalf("record event: %v", err)
	}

	events, err := svc.List(context.Background(), "", 10)
	if err != nil {
		t.Fatalf("list events: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if got := events[0].Type; got != "login" {
		t.Fatalf("expected login event, got %s", got)
	}
	_ = time.Now
}
