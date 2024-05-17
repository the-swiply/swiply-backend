package service

import (
	"context"

	"github.com/baobabus/go-apns/apns2"
	"github.com/google/uuid"

	"github.com/the-swiply/swiply-backend/notification/internal/domain"
)

type NotificationRepository interface {
	Create(ctx context.Context, notification domain.Notification) error
	Get(ctx context.Context, userID uuid.UUID) (domain.Notification, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}

type NotificationService struct {
	config NotificationServiceConfig
	repo   NotificationRepository
	apns   *apns2.Client
}

func NewNotificationService(config NotificationServiceConfig, repo NotificationRepository, client *apns2.Client) *NotificationService {
	return &NotificationService{config: config, repo: repo, apns: client}
}

func (n *NotificationService) Subscribe(ctx context.Context, notification domain.Notification) error {
	return n.repo.Create(ctx, notification)
}

func (n *NotificationService) Unsubscribe(ctx context.Context, userID uuid.UUID) error {
	return n.repo.Delete(ctx, userID)
}

func (n *NotificationService) Send(ctx context.Context, userID uuid.UUID, content string) error {
	notification, err := n.repo.Get(ctx, userID)
	if err != nil {
		return err
	}

	notif := &apns2.Notification{
		Recipient: notification.DeviceToken,
		Header:    &apns2.Header{Topic: n.config.Topic, Priority: apns2.PriorityHigh},
		Payload:   &apns2.Payload{APS: &apns2.APS{Alert: content}},
	}

	return n.apns.Push(notif, apns2.DefaultSigner, ctx, apns2.DefaultCallback)
}
