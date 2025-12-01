package ports

import (
	"context"

	"github.com/example/team-stack/backend/internal/app/core/model"
)

type UserService interface {
	List(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, u *model.User) error
	Login(ctx context.Context, email, password string) (*model.User, string, error)
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context) ([]model.User, error)
	Create(ctx context.Context, u *model.User) error
}

type JWTManager interface {
	Generate(userID, role string) (string, error)
}

type EventBus interface {
	Publish(subject string, data []byte) error
}

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Del(key string) error
}
