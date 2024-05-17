package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/multierr"

	"github.com/the-swiply/swiply-backend/chat/internal/cache"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
	"github.com/the-swiply/swiply-backend/chat/internal/glsync"
	"github.com/the-swiply/swiply-backend/chat/internal/pubsub"
	"github.com/the-swiply/swiply-backend/chat/internal/repository"
	"github.com/the-swiply/swiply-backend/chat/internal/server"
	"github.com/the-swiply/swiply-backend/chat/internal/service"
	"github.com/the-swiply/swiply-backend/chat/internal/sevents"
	"github.com/the-swiply/swiply-backend/chat/internal/workerpool"
	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/pkg/houston/runner"
)

const (
	authConfigPath = "configs/authorization.yaml"

	chatBroadcastChannel = "chat_broadcast"
)

type App struct {
	runner.RunStopper
	cfg *Config

	grpcServer *server.GRPCServer
	httpServer *server.HTTPServer

	redisSequenceCache      *cache.RedisSequenceCache
	redisMessagesPublisher  *pubsub.RedisMessagesPublisher
	redisMessagesSubscriber *pubsub.RedisMessagesSubscriber
	redisSyncer             *glsync.RedisSyncer
	db                      *pgxpool.Pool
	broadcastWP             *workerpool.Pool[domain.Message, error]

	stopCh chan struct{}
}

func NewApp(config *Config, runStopperPreset runner.RunStopper) *App {
	return &App{
		RunStopper: runStopperPreset,
		cfg:        config,
		stopCh:     make(chan struct{}),
	}
}

func (a *App) Run(ctx context.Context) error {
	defer close(a.stopCh)

	err := a.RunStopper.Run(ctx)
	if err != nil {
		return err
	}

	server.SetUserJWTSecret(os.Getenv("JWT_SECRET"))
	server.SetS2SJWTSecret(os.Getenv("S2S_JWT_SECRET"))
	err = server.ParseAuthorizationConfig(authConfigPath)
	if err != nil {
		return fmt.Errorf("can't parse auth config: %w", err)
	}

	db, err := dobby.NewPGXPool(ctx, dobby.PGXConfig{
		Username: a.cfg.Postgres.Username,
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     a.cfg.Postgres.Host,
		Port:     a.cfg.Postgres.Port,
		DBName:   a.cfg.Postgres.DBName,
		SSLMode:  a.cfg.Postgres.SSLMode,
	})
	if err != nil {
		return fmt.Errorf("can't init db: %w", err)
	}
	a.db = db

	chatRepo := repository.NewChatRepository(a.db)

	if os.Getenv("DISABLE_AUTO_MIGRATION") == "" {
		loggy.Infoln("starting migrations")
		err = dobby.AutoMigratePostgres(a.db, a.cfg.Postgres.MigrationsFolder)
		if err != nil {
			return fmt.Errorf("can't apply migrations: %w", err)
		}
		loggy.Infoln("migration done")
	}

	rdbSequence, err := cache.NewRedisSequenceCache(ctx, cache.RedisSequenceConfig{
		Addr:          a.cfg.Redis.Addr,
		Password:      os.Getenv("REDIS_PASSWORD"),
		DB:            int(a.cfg.Redis.DB.Sequence),
		SkipTLSVerify: a.cfg.Redis.SkipTLSVerify,
		Secure:        a.cfg.Redis.Secure,
	})
	if err != nil {
		return fmt.Errorf("can't init redis cache: %w", err)
	}
	a.redisSequenceCache = rdbSequence

	rdbPub, err := pubsub.NewRedisMessagesPublisher(ctx, pubsub.RedisPubSubConfig{
		Addr:          a.cfg.Redis.Addr,
		Password:      os.Getenv("REDIS_PASSWORD"),
		DB:            int(a.cfg.Redis.DB.MessagesPubSub),
		SkipTLSVerify: a.cfg.Redis.SkipTLSVerify,
		Secure:        a.cfg.Redis.Secure,
		ChannelName:   chatBroadcastChannel,
	})
	if err != nil {
		return fmt.Errorf("can't init redis publisher: %w", err)
	}
	a.redisMessagesPublisher = rdbPub
	rdbSync, err := glsync.NewRedisSyncer(ctx, glsync.RedisSyncerConfig{
		Addr:           a.cfg.Redis.Addr,
		Password:       os.Getenv("REDIS_PASSWORD"),
		DB:             int(a.cfg.Redis.DB.Syncer),
		SkipTLSVerify:  a.cfg.Redis.SkipTLSVerify,
		Secure:         a.cfg.Redis.Secure,
		LockExpiration: time.Millisecond * time.Duration(a.cfg.App.ChatLockExpirationMilliseconds),
	})
	if err != nil {
		return fmt.Errorf("can't init redis syncer: %w", err)
	}
	a.redisSyncer = rdbSync

	chatSvc := service.NewChatService(service.ChatConfig{}, a.redisSequenceCache, chatRepo, a.redisMessagesPublisher, a.redisSyncer.NewChatLock)

	a.broadcastWP = workerpool.NewPool[domain.Message, error](int(a.cfg.App.NumOfMessageSenderWorkers),
		chatSvc.SendMessageToChat,
		func(msg domain.Message) int64 {
			return msg.ChatID
		},
		workerpool.WithIgnoreResult[domain.Message, error](),
	)
	a.broadcastWP.Start()

	rdbSub, err := pubsub.NewRedisMessagesSubscriber(ctx, pubsub.RedisPubSubConfig{
		Addr:          a.cfg.Redis.Addr,
		Password:      os.Getenv("REDIS_PASSWORD"),
		DB:            int(a.cfg.Redis.DB.MessagesPubSub),
		SkipTLSVerify: a.cfg.Redis.SkipTLSVerify,
		Secure:        a.cfg.Redis.Secure,
		ChannelName:   chatBroadcastChannel,
	}, chatSvc, a.broadcastWP)
	if err != nil {
		return fmt.Errorf("can't init redis subscriber: %w", err)
	}
	a.redisMessagesSubscriber = rdbSub

	err = a.redisMessagesSubscriber.SubscribeOnMessages(ctx)
	if err != nil {
		return fmt.Errorf("subscriber can't subscribe on messages: %w", err)
	}

	errCh := make(chan error, 2)

	go func() {
		if err = a.runGRPCServer(chatSvc); err != nil {
			errCh <- fmt.Errorf("can't run grpc server: %w", err)
		} else {
			errCh <- nil
		}
	}()

	go func() {
		err = a.runHTTPServer(ctx, chatSvc)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("can't run http server: %w", err)
		} else {
			errCh <- nil
		}
	}()

	for i := 0; i < cap(errCh); i++ {
		if err = <-errCh; err != nil {
			return err
		}
	}

	return nil
}

func (a *App) runGRPCServer(chatSvc *service.ChatService) error {
	srv := server.NewGRPCServer(chatSvc)
	a.grpcServer = srv

	listener, err := net.Listen("tcp", a.cfg.GRPC.Addr)
	if err != nil {
		return fmt.Errorf("can't listen: %w", err)
	}
	defer listener.Close()

	loggy.Infoln("starting grpc server on", a.cfg.GRPC.Addr)

	if err = a.grpcServer.Serve(listener); err != nil {
		return fmt.Errorf("can't serve: %w", err)
	}

	return nil
}

func (a *App) runHTTPServer(ctx context.Context, chatService *service.ChatService) error {
	ws := sevents.NewWS(chatService)
	srv, err := server.NewHTTPServer(ctx, server.HTTPConfig{
		ServeAddr:    a.cfg.HTTP.Addr,
		GRPCEndpoint: a.cfg.GRPC.Addr,
		SwaggerPath:  a.cfg.Swagger.Path,
	}, ws)
	if err != nil {
		return err
	}

	a.httpServer = srv

	loggy.Infoln("starting http server on", a.cfg.HTTP.Addr)
	if err = a.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("can't serve: %w", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	err := multierr.Combine(
		a.grpcServer.Shutdown(ctx),
		a.httpServer.Shutdown(ctx),
		a.redisSyncer.Stop(ctx),
		a.redisSequenceCache.Stop(ctx),
		a.redisMessagesPublisher.Stop(ctx),
		a.redisMessagesSubscriber.Stop(ctx),
		a.broadcastWP.Stop(ctx),
		a.RunStopper.Stop(ctx),
	)
	a.db.Close()

	if err != nil {
		return err
	}

	select {
	case <-a.stopCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
