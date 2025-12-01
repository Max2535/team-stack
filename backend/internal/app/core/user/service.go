package user

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/example/team-stack/backend/internal/app/core/model"
	"github.com/example/team-stack/backend/internal/app/ports"
)

var (
	ErrInvalidLogin = errors.New("invalid email or password")
)

type Service struct {
	repo ports.UserRepository
	jwt  ports.JWTManager
	bus  ports.EventBus
}

func NewService(repo ports.UserRepository, jwt ports.JWTManager, bus ports.EventBus) *Service {
	return &Service{repo: repo, jwt: jwt, bus: bus}
}

func (s *Service) List(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, u *model.User) error {
	if err := s.repo.Create(ctx, u); err != nil {
		return err
	}
	ev := map[string]any{
		"type":  "UserCreated",
		"id":    u.ID,
		"email": u.Email,
	}
	b, _ := json.Marshal(ev)
	_ = s.bus.Publish("user.created", b)
	return nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	u, err := s.repo.FindByEmail(ctx, email)
	if err != nil || u == nil {
		return nil, "", ErrInvalidLogin
	}
	if !secureComparePassword(u.PasswordHash, password) {
		return nil, "", ErrInvalidLogin
	}
	token, err := s.jwt.Generate(u.ID, u.Role)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

func secureComparePassword(storedHash, password string) bool {
	if storedHash == "" || password == "" {
		return false
	}

	expected, err := hex.DecodeString(storedHash)
	if err != nil {
		return false
	}

	derived := sha256.Sum256([]byte(password))
	if len(expected) != len(derived) {
		return false
	}

	return subtle.ConstantTimeCompare(expected, derived[:]) == 1
}
