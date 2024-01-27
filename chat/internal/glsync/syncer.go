package glsync

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/redis/go-redis/v9"
	"github.com/the-swiply/swiply-backend/chat/internal/service"
	"strconv"
)

type RedisSyncer struct {
	client *redis.Client
	sync   *redsync.Redsync
	cfg    RedisSyncerConfig
}

func NewRedisSyncer(ctx context.Context, cfg RedisSyncerConfig) (*RedisSyncer, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("can't ping redis: %w", err)
	}

	rs := redsync.New(goredis.NewPool(rc))

	return &RedisSyncer{
		client: rc,
		sync:   rs,
		cfg:    cfg,
	}, nil
}

func (r *RedisSyncer) NewChatLock(chatID int64) service.ChatLocker {
	return &ChatLock{Mutex: r.sync.NewMutex(strconv.FormatInt(chatID, 10), redsync.WithExpiry(r.cfg.LockExpiration))}
}

func (r *RedisSyncer) Stop(_ context.Context) error {
	return r.client.Close()
}
