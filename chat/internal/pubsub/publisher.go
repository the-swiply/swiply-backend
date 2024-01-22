package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
)

type RedisMessagesPublisher struct {
	cfg    RedisPubSubConfig
	client *redis.Client
}

func NewRedisMessagesPublisher(ctx context.Context, cfg RedisPubSubConfig) (*RedisMessagesPublisher, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("can't ping redis: %w", err)
	}

	return &RedisMessagesPublisher{
		client: rc,
		cfg:    cfg,
	}, nil
}

func (r *RedisMessagesPublisher) PublishMessage(ctx context.Context, msg domain.ChatMessage) error {
	marshalledMsg, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("can't marshal message: %w", err)
	}

	return r.client.Publish(ctx, r.cfg.ChannelName, marshalledMsg).Err()
}

func (r *RedisMessagesPublisher) Stop(_ context.Context) error {
	return r.client.Close()
}
