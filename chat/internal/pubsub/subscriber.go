package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
	"github.com/the-swiply/swiply-backend/chat/internal/service"
	"github.com/the-swiply/swiply-backend/chat/internal/workerpool"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"time"
)

const (
	defaultSendTimeout = time.Minute * 3
)

type RedisMessagesSubscriber struct {
	cfg         RedisPubSubConfig
	client      *redis.Client
	redisPubSub *redis.PubSub
	workerPool  *workerpool.Pool[domain.Message, error]

	chatService *service.ChatService
	stopCh      chan struct{}
}

func NewRedisMessagesSubscriber(ctx context.Context, cfg RedisPubSubConfig, chatService *service.ChatService,
	workerPool *workerpool.Pool[domain.Message, error]) (*RedisMessagesSubscriber, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := rc.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("can't ping redis: %w", err)
	}

	return &RedisMessagesSubscriber{
		client:      rc,
		cfg:         cfg,
		workerPool:  workerPool,
		chatService: chatService,
		stopCh:      make(chan struct{}),
	}, nil
}

func (r *RedisMessagesSubscriber) SubscribeOnMessages(ctx context.Context) error {
	sub := r.client.Subscribe(ctx, r.cfg.ChannelName)
	_, err := sub.Receive(ctx)
	if err != nil {
		return fmt.Errorf("can't subscribe on pubsub channel")
	}

	r.redisPubSub = sub

	redisMsgCh := sub.Channel(redis.WithChannelSendTimeout(defaultSendTimeout))

	go func() {
		var msg domain.Message
		for redisMsg := range redisMsgCh {
			err = json.Unmarshal([]byte(redisMsg.Payload), &msg)
			if err != nil {
				loggy.Errorf("can't unmarshal message: %v", err)
				continue
			}

			err = r.workerPool.AddTask(ctx, msg)
			if err != nil {
				loggy.Errorf("can't add task to pool: %v", err)
				continue
			}
		}

		close(r.stopCh)
	}()

	return nil
}

func (r *RedisMessagesSubscriber) Stop(ctx context.Context) error {
	err := r.redisPubSub.Close()
	if err != nil {
		return fmt.Errorf("can't unsubscribe from messages: %w", err)
	}

	select {
	case <-r.stopCh:
	case <-ctx.Done():
		return ctx.Err()
	}

	return r.client.Close()
}
