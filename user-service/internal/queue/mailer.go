package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/the-swiply/swiply-backend/user/internal/entity"
	"github.com/the-swiply/swiply-backend/user/internal/service"
	"time"
)

const (
	typeEmailSend = "email:send"

	emailSendRetention = time.Hour * 24
)

type MailerQueue struct {
	cfg    MailerConfig
	server *asynq.Server
	client *asynq.Client
	mux    *asynq.ServeMux
}

func NewMailerQueue(cfg MailerConfig, senderService *service.SenderService) *MailerQueue {
	redisConnOpts := asynq.RedisClientOpt{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	}
	server := asynq.NewServer(
		redisConnOpts,
		asynq.Config{
			LogLevel:        asynq.ErrorLevel,
			Concurrency:     1,
			ShutdownTimeout: cfg.SendTimeout,
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle(typeEmailSend,
		&emailSendHandler{
			senderService:  senderService,
			pauseAfterSend: cfg.AfterSendWorkerPause,
		},
	)

	client := asynq.NewClient(redisConnOpts)

	return &MailerQueue{
		server: server,
		client: client,
		cfg:    cfg,
		mux:    mux,
	}
}

func (m *MailerQueue) Run() error {
	return m.server.Start(m.mux)
}

func (m *MailerQueue) Stop(ctx context.Context) error {
	stopCh := make(chan struct{})

	go func() {
		m.server.Shutdown()
		close(stopCh)
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-stopCh:
	}

	return m.client.Close()
}

func (m *MailerQueue) ScheduleEmailSend(ctx context.Context, info entity.SendAuthCodeInfo) error {
	payload, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("can't marshal payload: %w", err)
	}

	task := asynq.NewTask(typeEmailSend, payload, asynq.MaxRetry(3), asynq.Timeout(m.cfg.SendTimeout), asynq.Retention(emailSendRetention))
	if err != nil {
		return fmt.Errorf("can't create task: %w", err)
	}

	_, err = m.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("can't enqueue task: %w", err)
	}

	return nil
}

type emailSendHandler struct {
	senderService  *service.SenderService
	pauseAfterSend time.Duration
}

func (e *emailSendHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	var info entity.SendAuthCodeInfo
	if err := json.Unmarshal(t.Payload(), &info); err != nil {
		return fmt.Errorf("can't unmarshal payload: %v: %w", err, asynq.SkipRetry)
	}

	err := e.senderService.SendEmailWithAuthorizationCode(ctx, info.To, info.Subject, info.Code)
	if err != nil {
		return fmt.Errorf("can't send email: %w", err)
	}

	time.Sleep(e.pauseAfterSend)

	return nil
}
