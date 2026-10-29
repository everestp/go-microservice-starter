package domain

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrUserExists           = errors.New("user already exists")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrRefreshTokenRevoked  = errors.New("refresh token revoked")
	ErrInvalidRefreshToken  = errors.New("invalid refresh token")
)

// User represents an authenticated account.
type User struct {
	ID           string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

// RefreshToken represents an opaque refresh token issued to a user.
type RefreshToken struct {
	ID        string
	UserID    string
	FamilyID  string
	TokenHash string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

// TokenPair contains the issued bearer tokens for a user session.
type TokenPair struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}
