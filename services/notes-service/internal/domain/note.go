package domain

import (
	"errors"
	"time"
)

var (
	ErrNoteNotFound        = errors.New("note not found")
	ErrNoteVersionMismatch = errors.New("note version mismatch")
)

// Note is a single user-owned document in the notes service.
type Note struct {
	ID        string
	Title     string
	Body      string
	Version   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
