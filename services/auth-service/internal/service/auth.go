package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/example/go-microservice-starter/pkg/clock"
	"github.com/example/go-microservice-starter/services/auth-service/internal/domain"
	"github.com/google/uuid"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 7 * 24 * time.Hour
)

// Repository persists users and refresh tokens.
type Repository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetUserByID(ctx context.Context, id string) (domain.User, error)
	SaveRefreshToken(ctx context.Context, token domain.RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, hash string) (domain.RefreshToken, error)
	RevokeRefreshFamily(ctx context.Context, familyID string) error
}

// PasswordHasher hashes and verifies user passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

// TokenIssuer issues access and refresh tokens.
type TokenIssuer interface {
	IssueAccess(userID string, now time.Time) (string, error)
	IssueRefresh(userID, familyID string, now time.Time) (string, error)
}

// AuthService implements the authentication domain rules.
type AuthService struct {
	repo   Repository
	hasher PasswordHasher
	issuer TokenIssuer
	clock  clock.Clock
}

// NewAuthService constructs an AuthService with repository, password hashing, token issuance, and time abstraction.
func NewAuthService(repo Repository, hasher PasswordHasher, issuer TokenIssuer, c clock.Clock) *AuthService {
	return &AuthService{repo: repo, hasher: hasher, issuer: issuer, clock: c}
}

// Register creates a new user and issues a fresh token pair.
func (s *AuthService) Register(ctx context.Context, email, password string) (domain.TokenPair, error) {
	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return domain.TokenPair{}, domain.ErrUserExists
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return domain.TokenPair{}, err
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return domain.TokenPair{}, err
	}

	user := domain.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    s.clock.Now(),
	}
	if err := s.repo.CreateUser(ctx, user); err != nil {
		return domain.TokenPair{}, err
	}

	access, err := s.issuer.IssueAccess(user.ID, s.clock.Now())
	if err != nil {
		return domain.TokenPair{}, err
	}
	familyID := uuid.NewString()
	refresh, err := s.issuer.IssueRefresh(user.ID, familyID, s.clock.Now())
	if err != nil {
		return domain.TokenPair{}, err
	}

	if err := s.repo.SaveRefreshToken(ctx, domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		FamilyID:  familyID,
		TokenHash: hashToken(refresh),
		ExpiresAt: s.clock.Now().Add(RefreshTokenTTL),
		CreatedAt: s.clock.Now(),
	}); err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{UserID: user.ID, AccessToken: access, RefreshToken: refresh}, nil
}

// Login authenticates a user and returns a fresh token pair.
func (s *AuthService) Login(ctx context.Context, email, password string) (domain.TokenPair, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.TokenPair{}, domain.ErrInvalidCredentials
	}
	if !s.hasher.Compare(user.PasswordHash, password) {
		return domain.TokenPair{}, domain.ErrInvalidCredentials
	}

	access, err := s.issuer.IssueAccess(user.ID, s.clock.Now())
	if err != nil {
		return domain.TokenPair{}, err
	}
	familyID := uuid.NewString()
	refresh, err := s.issuer.IssueRefresh(user.ID, familyID, s.clock.Now())
	if err != nil {
		return domain.TokenPair{}, err
	}

	if err := s.repo.SaveRefreshToken(ctx, domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		FamilyID:  familyID,
		TokenHash: hashToken(refresh),
		ExpiresAt: s.clock.Now().Add(RefreshTokenTTL),
		CreatedAt: s.clock.Now(),
	}); err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{UserID: user.ID, AccessToken: access, RefreshToken: refresh}, nil
}

// Refresh rotates a refresh token for the same user and family.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (domain.TokenPair, error) {
	stored, err := s.repo.GetRefreshTokenByHash(ctx, hashToken(refreshToken))
	if err != nil {
		return domain.TokenPair{}, domain.ErrInvalidRefreshToken
	}
	if stored.Revoked {
		return domain.TokenPair{}, domain.ErrRefreshTokenRevoked
	}
	if s.clock.Now().After(stored.ExpiresAt) {
		return domain.TokenPair{}, domain.ErrInvalidRefreshToken
	}

	user, err := s.repo.GetUserByID(ctx, stored.UserID)
	if err != nil {
		return domain.TokenPair{}, domain.ErrUserNotFound
	}

	newFamilyID := stored.FamilyID
	newRefresh, err := s.issuer.IssueRefresh(user.ID, newFamilyID, s.clock.Now())
	if err != nil {
		return domain.TokenPair{}, err
	}
	if err := s.repo.SaveRefreshToken(ctx, domain.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    user.ID,
		FamilyID:  newFamilyID,
		TokenHash: hashToken(newRefresh),
		ExpiresAt: s.clock.Now().Add(RefreshTokenTTL),
		CreatedAt: s.clock.Now(),
	}); err != nil {
		return domain.TokenPair{}, err
	}

	access, err := s.issuer.IssueAccess(user.ID, s.clock.Now())
	if err != nil {
		return domain.TokenPair{}, err
	}

	if err := s.repo.RevokeRefreshFamily(ctx, stored.FamilyID); err != nil {
		return domain.TokenPair{}, err
	}

	return domain.TokenPair{UserID: user.ID, AccessToken: access, RefreshToken: newRefresh}, nil
}

// Logout invalidates the given refresh token and its family.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	stored, err := s.repo.GetRefreshTokenByHash(ctx, hashToken(refreshToken))
	if err != nil {
		return domain.ErrRefreshTokenNotFound
	}
	stored.Revoked = true
	if err := s.repo.RevokeRefreshFamily(ctx, stored.FamilyID); err != nil {
		return err
	}
	return s.repo.SaveRefreshToken(ctx, stored)
}

// Me returns the current user profile by ID.
func (s *AuthService) Me(ctx context.Context, userID string) (domain.User, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func hashToken(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
