package postgres

import (
	"context"
	"sync"

	"github.com/example/go-microservice-starter/services/auth-service/internal/domain"
)

// Repository is an in-memory persistence implementation used by the auth-service scaffold.
type Repository struct {
	mu      sync.RWMutex
	users   map[string]domain.User
	byEmail map[string]domain.User
	tokens  map[string]domain.RefreshToken
}

// NewRepository constructs an in-memory repository for local development and tests.
func NewRepository() *Repository {
	return &Repository{
		users:   map[string]domain.User{},
		byEmail: map[string]domain.User{},
		tokens:  map[string]domain.RefreshToken{},
	}
}

func (r *Repository) CreateUser(_ context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.byEmail[user.Email]; exists {
		return domain.ErrUserExists
	}
	r.users[user.ID] = user
	r.byEmail[user.Email] = user
	return nil
}

func (r *Repository) GetUserByEmail(_ context.Context, email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.byEmail[email]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *Repository) GetUserByID(_ context.Context, id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, ok := r.users[id]
	if !ok {
		return domain.User{}, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *Repository) SaveRefreshToken(_ context.Context, token domain.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tokens[token.TokenHash] = token
	return nil
}

func (r *Repository) GetRefreshTokenByHash(_ context.Context, hash string) (domain.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	token, ok := r.tokens[hash]
	if !ok {
		return domain.RefreshToken{}, domain.ErrRefreshTokenNotFound
	}
	return token, nil
}

func (r *Repository) RevokeRefreshFamily(_ context.Context, familyID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for key, token := range r.tokens {
		if token.FamilyID == familyID {
			token.Revoked = true
			r.tokens[key] = token
		}
	}
	return nil
}
