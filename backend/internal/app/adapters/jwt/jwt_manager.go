package jwt

import (
    "time"

    "github.com/example/team-stack/backend/internal/config"
    "github.com/example/team-stack/backend/internal/app/ports"
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
