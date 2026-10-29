package service

import (
	"context"
	"testing"
	"time"

	"github.com/example/go-microservice-starter/pkg/clock"
	"github.com/example/go-microservice-starter/services/auth-service/internal/domain"
)

type fakeRepo struct {
	users   map[string]domain.User
	refresh map[string]domain.RefreshToken
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{users: map[string]domain.User{}, refresh: map[string]domain.RefreshToken{}}
}

func (r *fakeRepo) CreateUser(_ context.Context, u domain.User) error {
	r.users[u.Email] = u
	return nil
}

func (r *fakeRepo) GetUserByEmail(_ context.Context, email string) (domain.User, error) {
	if u, ok := r.users[email]; ok {
		return u, nil
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (r *fakeRepo) GetUserByID(_ context.Context, id string) (domain.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (r *fakeRepo) SaveRefreshToken(_ context.Context, tok domain.RefreshToken) error {
	r.refresh[tok.TokenHash] = tok
	return nil
}

func (r *fakeRepo) GetRefreshTokenByHash(_ context.Context, hash string) (domain.RefreshToken, error) {
	if tok, ok := r.refresh[hash]; ok {
		return tok, nil
	}
	return domain.RefreshToken{}, domain.ErrRefreshTokenNotFound
}

func (r *fakeRepo) RevokeRefreshFamily(_ context.Context, familyID string) error {
	for key, tok := range r.refresh {
		if tok.FamilyID == familyID {
			tok.Revoked = true
			r.refresh[key] = tok
		}
	}
	return nil
}

type fakeHasher struct{}

func (fakeHasher) Hash(password string) (string, error) {
	return "hash:" + password, nil
}

func (fakeHasher) Compare(hash, password string) bool {
	return hash == "hash:"+password
}

type fakeIssuer struct{}

func (fakeIssuer) IssueAccess(userID string, _ time.Time) (string, error) {
	return "access-" + userID, nil
}

func (fakeIssuer) IssueRefresh(userID, familyID string, _ time.Time) (string, error) {
	return "refresh-" + userID + "-" + familyID, nil
}

func TestRegisterAndLogin(t *testing.T) {
	repo := newFakeRepo()
	svc := NewAuthService(repo, fakeHasher{}, fakeIssuer{}, clock.RealClock{})

	pair, err := svc.Register(context.Background(), "alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if pair.AccessToken == "" || pair.RefreshToken == "" {
		t.Fatal("tokens must be non-empty")
	}

	loggedIn, err := svc.Login(context.Background(), "alice@example.com", "secret123")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if loggedIn.UserID != pair.UserID {
		t.Fatalf("expected user ID %s, got %s", pair.UserID, loggedIn.UserID)
	}
}

func TestLoginRejectsBadPassword(t *testing.T) {
	repo := newFakeRepo()
	svc := NewAuthService(repo, fakeHasher{}, fakeIssuer{}, clock.RealClock{})

	_, err := svc.Register(context.Background(), "bob@example.com", "correct")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, err = svc.Login(context.Background(), "bob@example.com", "wrong")
	if err == nil {
		t.Fatal("expected invalid credentials error")
	}
}
