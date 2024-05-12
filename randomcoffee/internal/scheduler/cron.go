package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/service"
)

const (
	typeRandomCoffeeTrigger    = "random:coffee:trigger"
	randomCoffeeTriggerTimeout = time.Minute * 5

	defaultRetention = time.Hour * 24
)

type RedisCron struct {
	cfg                 RedisCronConfig
	scheduler           *asynq.Scheduler
	server              *asynq.Server
	randomCoffeeService *service.RandomCoffeeService
}

func NewRedisCron(cfg RedisCronConfig, randomCoffeeService *service.RandomCoffeeService) (*RedisCron, error) {
	redisConnOpts := asynq.RedisClientOpt{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	server := asynq.NewServer(
		redisConnOpts,
		asynq.Config{
			LogLevel:        asynq.ErrorLevel,
			Concurrency:     1,
			ShutdownTimeout: 5 * time.Second,
		},
	)

	scheduler := asynq.NewScheduler(
		redisConnOpts,
		&asynq.SchedulerOpts{
			Location: time.UTC,
			LogLevel: asynq.ErrorLevel,
		},
	)

	return &RedisCron{
		cfg:                 cfg,
		scheduler:           scheduler,
		server:              server,
		randomCoffeeService: randomCoffeeService,
	}, nil
}

func (r *RedisCron) RunScheduler() error {
	_, err := r.scheduler.Register(r.cfg.RandomCoffeeTriggerCron, asynq.NewTask(typeRandomCoffeeTrigger, nil,
		asynq.Unique(randomCoffeeTriggerTimeout),
		asynq.MaxRetry(2),
		asynq.Timeout(randomCoffeeTriggerTimeout),
		asynq.Retention(defaultRetention)))
	if err != nil {
		return fmt.Errorf("can't register random coffee trigger task: %w", err)
	}

	return r.scheduler.Start()
}

func (r *RedisCron) RunServer() error {
	mux := asynq.NewServeMux()
	mux.Handle(typeRandomCoffeeTrigger, &randomCoffeeTriggerHandler{
		randomCoffeeService: r.randomCoffeeService,
	})

	return r.server.Start(mux)
}

func (r *RedisCron) Stop(ctx context.Context) error {
	stopCh := make(chan struct{})

	go func() {
		r.scheduler.Shutdown()
		r.server.Shutdown()
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCh:
	}

	return nil
}

type randomCoffeeTriggerHandler struct {
	randomCoffeeService *service.RandomCoffeeService
}

func (r *randomCoffeeTriggerHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	loggy.Infoln("start random coffee task")
	err := r.randomCoffeeService.Schedule(ctx)
	if err != nil {
		err := fmt.Errorf("can't update statistics: %w", err)
		loggy.Errorln(err)
		return err
	}
	loggy.Infoln("random coffee task completed")

	return nil
}
