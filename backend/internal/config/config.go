package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName string
	Env     string
	Host    string
	Port    string

	DB struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		ConnMaxLife  time.Duration
	}

	RedisAddr string

	KafkaBrokers []string
	KafkaTopic   string

	JWT struct {
		Secret    string
		AccessTTL time.Duration
	}

	LogLevel string
}

func Load() (*Config, error) {
	cfg := &Config{
		AppName:      getenv("APP_NAME", "team-api"),
		Env:          getenv("APP_ENV", "dev"),
		Host:         getenv("APP_HOST", "0.0.0.0"),
		Port:         getenv("APP_PORT", "8080"),
		RedisAddr:    getenv("REDIS_ADDR", "localhost:6379"),
		KafkaBrokers: []string{getenv("KAFKA_BROKER", "localhost:9092")},
		KafkaTopic:   getenv("KAFKA_TOPIC", "team-events"),
		LogLevel:     getenv("LOG_LEVEL", "debug"),
	}

	cfg.DB.DSN = getenv("DB_DSN", "postgres://team:team@localhost:5432/teamdb?sslmode=disable")
	cfg.DB.MaxOpenConns = getenvInt("DB_MAX_OPEN", 20)
	cfg.DB.MaxIdleConns = getenvInt("DB_MAX_IDLE", 5)
	cfg.DB.ConnMaxLife = getenvDuration("DB_CONN_MAX_LIFE", time.Minute*30)

	cfg.JWT.Secret = getenv("JWT_SECRET", "change-me")
	cfg.JWT.AccessTTL = getenvDuration("JWT_ACCESS_TTL", time.Minute*15)

	return cfg, nil
}

func (c *Config) Addr() string {
	return c.Host + ":" + c.Port
}

func getenv(key, def string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return def
}

func getenvInt(key string, def int) int {
	if v, ok := os.LookupEnv(key); ok {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return def
}

func getenvDuration(key string, def time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok {
		d, err := time.ParseDuration(v)
		if err == nil {
			return d
		}
	}
	return def
}
