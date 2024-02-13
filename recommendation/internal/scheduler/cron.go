package scheduler

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/the-swiply/swiply-backend/pkg/houston/loggy"
	"github.com/the-swiply/swiply-backend/recommendation/internal/service"
	"time"
)

const (
	typeStatisticUpdate    = "statistic:update"
	statisticUpdateTimeout = time.Minute * 5

	typeTriggerOracleLearn    = "oracle:learn:trigger"
	triggerOracleLearnTimeout = time.Second * 5

	defaultRetention = time.Hour * 24
)

type RedisCron struct {
	cfg       RedisCronConfig
	scheduler *asynq.Scheduler
	server    *asynq.Server

	dpService *service.DataProviderService
}

func NewRedisCron(cfg RedisCronConfig, dpService *service.DataProviderService) (*RedisCron, error) {
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
		cfg:       cfg,
		scheduler: scheduler,
		server:    server,
		dpService: dpService,
	}, nil
}

func (r *RedisCron) RunScheduler() error {
	_, err := r.scheduler.Register(r.cfg.StatisticUpdateCron, asynq.NewTask(typeStatisticUpdate, nil,
		asynq.Unique(statisticUpdateTimeout),
		asynq.MaxRetry(2),
		asynq.Timeout(statisticUpdateTimeout),
		asynq.Retention(defaultRetention)))
	if err != nil {
		return fmt.Errorf("can't register statistic update task: %w", err)
	}

	_, err = r.scheduler.Register(r.cfg.StatisticUpdateCron, asynq.NewTask(typeTriggerOracleLearn, nil,
		asynq.Unique(triggerOracleLearnTimeout),
		asynq.MaxRetry(2),
		asynq.Timeout(triggerOracleLearnTimeout),
		asynq.Retention(defaultRetention)))
	if err != nil {
		return fmt.Errorf("can't register learn oracle update task: %w", err)
	}

	return r.scheduler.Start()
}

func (r *RedisCron) RunServer() error {
	mux := asynq.NewServeMux()
	mux.Handle(typeStatisticUpdate, &statisticUpdateHandler{
		dpService: r.dpService,
	})
	mux.Handle(typeTriggerOracleLearn, &triggerOracleLearnHandler{
		dpService: r.dpService,
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

type statisticUpdateHandler struct {
	dpService *service.DataProviderService
}

func (s *statisticUpdateHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	loggy.Infoln("start updating statistics")
	err := s.dpService.UpdateStatistic(ctx)
	if err != nil {
		return fmt.Errorf("can't update statistics: %w", err)
	}
	loggy.Infoln("statistics updated")

	return nil
}

type triggerOracleLearnHandler struct {
	dpService *service.DataProviderService
}

func (s *triggerOracleLearnHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {
	err := s.dpService.UpdateOracleData(ctx)
	if err != nil {
		return fmt.Errorf("can't update oracle data: %w", err)
	}

	return nil
}
