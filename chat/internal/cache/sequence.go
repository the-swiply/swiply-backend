package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisSequenceCache struct {
	client *redis.Client
	cfg    RedisSequenceConfig
}

func NewRedisSequenceCache(ctx context.Context, cfg RedisSequenceConfig) (*RedisSequenceCache, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("can't ping redis: %w", err)
	}

	return &RedisSequenceCache{
		client: rc,
		cfg:    cfg,
	}, nil
}

func (r *RedisSequenceCache) Stop(_ context.Context) error {
	return r.client.Close()
}
