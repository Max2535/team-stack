package db

import (
	"context"

	"github.com/example/team-stack/backend/internal/app/core/model"
	"github.com/example/team-stack/backend/internal/app/ports"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserPostgresRepo struct {
	db *sqlx.DB
}

func NewUserPostgresRepo(db *sqlx.DB) ports.UserRepository {
	return &UserPostgresRepo{db: db}
}

func (r *UserPostgresRepo) FindByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserPostgresRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := r.db.GetContext(ctx, &u, `SELECT * FROM users WHERE email = $1`, email)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserPostgresRepo) List(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.db.SelectContext(ctx, &users, `SELECT * FROM users ORDER BY created_at DESC`)
	return users, err
}

func (r *UserPostgresRepo) Create(ctx context.Context, u *model.User) error {
	if u.ID == "" {
		u.ID = uuid.NewString()
	}
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO users (id, email, password_hash, name, role, created_at, updated_at)
        VALUES ($1,$2,$3,$4,$5,now(),now())
    `, u.ID, u.Email, u.PasswordHash, u.Name, u.Role)
	return err
}
