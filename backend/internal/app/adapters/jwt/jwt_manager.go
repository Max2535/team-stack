package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/example/team-stack/backend/internal/app/ports"
	"github.com/example/team-stack/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTManager(cfg *config.Config) ports.JWTManager {
	return &Manager{
		secret: []byte(cfg.JWT.Secret),
		ttl:    cfg.JWT.AccessTTL,
	}
}

func (m *Manager) Generate(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(m.ttl).Unix(),
		"iat":  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) Verify(token string) (*ports.AuthClaims, error) {
	claims := jwt.MapClaims{}

	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	if err := claims.Valid(); err != nil {
		return nil, err
	}

	sub, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)
	if sub == "" || role == "" {
		return nil, errors.New("missing token claims")
	}

	return &ports.AuthClaims{UserID: sub, Role: role}, nil
}
