package user

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/example/team-stack/backend/internal/app/core/model"
	"github.com/example/team-stack/backend/internal/app/ports"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidLogin = errors.New("invalid email or password")
)

// BcryptCost is the cost parameter for bcrypt hashing
// Default is 10, increase for more security (slower)
const BcryptCost = 10

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

func (s *Service) GetByID(ctx context.Context, id string) (*model.User, error) {
	return s.repo.FindByID(ctx, id)
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
	if !CheckPassword(u.PasswordHash, password) {
		return nil, "", ErrInvalidLogin
	}
	token, err := s.jwt.Generate(u.ID, u.Role)
	if err != nil {
		return nil, "", err
	}
	return u, token, nil
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a bcrypt hashed password with its possible plaintext equivalent
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
