package cache

import (
	"context"

	"github.com/example/team-stack/backend/internal/app/ports"
	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	cli *redis.Client
}

func NewRedis(addr string) ports.Cache {
	return &RedisCache{
		cli: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}
}

func (r *RedisCache) Get(key string) ([]byte, error) {
	val, err := r.cli.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	return val, err
}

func (r *RedisCache) Set(key string, value []byte) error {
	return r.cli.Set(context.Background(), key, value, 0).Err()
}

func (r *RedisCache) Del(key string) error {
	return r.cli.Del(context.Background(), key).Err()
}
