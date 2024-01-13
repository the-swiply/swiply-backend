package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/the-swiply/swiply-backend/user/internal/domain"
	"strconv"
	"time"
)

type RedisCodeCache struct {
	client *redis.Client
	cfg    RedisCodesConfig
}

func NewRedisCodeCache(ctx context.Context, cfg RedisCodesConfig) (*RedisCodeCache, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("can't ping redis: %w", err)
	}

	return &RedisCodeCache{
		client: rc,
		cfg:    cfg,
	}, nil
}

func (r *RedisCodeCache) StoreAuthCode(ctx context.Context, email string, code int) error {
	return r.client.Set(ctx, email, strconv.Itoa(code), r.cfg.AuthCodeTTL).Err()
}

func (r *RedisCodeCache) GetAuthCode(ctx context.Context, email string) (string, error) {
	code, err := r.client.Get(ctx, email).Result()
	if errors.Is(err, redis.Nil) {
		return "", domain.ErrEntityIsNotExists
	}

	return code, err
}

func (r *RedisCodeCache) DeleteAuthCode(ctx context.Context, email string) error {
	return r.client.Del(ctx, email).Err()
}

func (r *RedisCodeCache) GetAuthCodeTTL(ctx context.Context, email string) (time.Duration, error) {
	ttl, err := r.client.TTL(ctx, email).Result()
	if err != nil {
		return 0, err
	}

	return ttl, err
}

func (r *RedisCodeCache) Stop(_ context.Context) error {
	return r.client.Close()
}
