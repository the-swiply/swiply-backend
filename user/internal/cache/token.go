package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/the-swiply/swiply-backend/user/internal/domain"
)

type RedisTokenCache struct {
	client *redis.Client
	cfg    RedisTokensConfig
}

func NewRedisTokenCache(ctx context.Context, cfg RedisTokensConfig) (*RedisTokenCache, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("can't ping redis: %w", err)
	}

	return &RedisTokenCache{
		client: rc,
		cfg:    cfg,
	}, nil
}

func (r *RedisTokenCache) StoreFingerprint(ctx context.Context, token string, fingerprint string) error {
	return r.client.Set(ctx, token, fingerprint, r.cfg.RefreshTokenTTL).Err()
}

func (r *RedisTokenCache) GetFingerprint(ctx context.Context, token string) (string, error) {
	fingerprint, err := r.client.Get(ctx, token).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrEntityIsNotExists
	}

	return fingerprint, err
}

func (r *RedisTokenCache) DeleteFingerprint(ctx context.Context, token string) error {
	return r.client.Del(ctx, token).Err()
}

func (r *RedisTokenCache) Stop(_ context.Context) error {
	return r.client.Close()
}
