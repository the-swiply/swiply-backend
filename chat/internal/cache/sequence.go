package cache

import (
	"context"
	"crypto/tls"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisSequenceCache struct {
	client *redis.Client
	cfg    RedisSequenceConfig
}

func NewRedisSequenceCache(ctx context.Context, cfg RedisSequenceConfig) (*RedisSequenceCache, error) {
	var tlsConf *tls.Config
	if cfg.Secure {
		tlsConf = &tls.Config{
			InsecureSkipVerify: cfg.SkipTLSVerify,
		}
	}

	rc := redis.NewClient(&redis.Options{
		Addr:      cfg.Addr,
		Password:  cfg.Password,
		DB:        cfg.DB,
		TLSConfig: tlsConf,
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

func (r *RedisSequenceCache) GenerateNextID(ctx context.Context, chatID int64) (int64, error) {
	nextID, err := r.client.Incr(ctx, strconv.FormatInt(chatID, 10)).Result()
	if err != nil {
		return 0, err
	}

	return nextID, nil
}

func (r *RedisSequenceCache) RollbackID(ctx context.Context, chatID int64) error {
	return r.client.Decr(ctx, strconv.FormatInt(chatID, 10)).Err()
}

func (r *RedisSequenceCache) Stop(_ context.Context) error {
	return r.client.Close()
}
