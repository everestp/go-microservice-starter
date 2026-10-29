package domain

import (
	"encoding/json"
	"errors"
	"time"
)

var ErrEventNotFound = errors.New("event not found")

type AuditEvent struct {
	ID        string
	Type      string
	Version   int32
	Source    string
	CreatedAt time.Time
	UserID    string
	Data      map[string]string
}

func (e AuditEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID        string            `json:"id"`
		Type      string            `json:"type"`
		Version   int32             `json:"version"`
		Source    string            `json:"source"`
		CreatedAt time.Time         `json:"created_at"`
		UserID    string            `json:"user_id"`
		Data      map[string]string `json:"data"`
	}{
		ID:        e.ID,
		Type:      e.Type,
		Version:   e.Version,
		Source:    e.Source,
		CreatedAt: e.CreatedAt,
		UserID:    e.UserID,
		Data:      e.Data,
	})
}
