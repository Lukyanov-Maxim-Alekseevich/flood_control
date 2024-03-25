package redis

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"task/internal/infra/storage"
	"time"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{client: client}
}

func (rs *RedisStorage) Set(ctx context.Context, key, value string, exp time.Duration) error {
	err := rs.client.Set(ctx, key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rs *RedisStorage) Get(ctx context.Context, key string) (string, error) {
	val, err := rs.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", storage.ErrRedisKeyNotFound
		}
		return "", err
	}
	return val, nil
}
