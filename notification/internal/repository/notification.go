package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"

	"github.com/the-swiply/swiply-backend/notification/internal/domain"
)

const (
	notificationTable = "notification"
)

type NotificationRepository struct {
	db *pgxpool.Pool
}

func NewNotificationRepository(db *pgxpool.Pool) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (n *NotificationRepository) executor(ctx context.Context) dobby.Executor {
	tx := dobby.ExtractPGXTx(ctx)
	if tx != nil {
		return tx
	}

	return n.db
}

func (n *NotificationRepository) Create(ctx context.Context, notification domain.Notification) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, device_token)
VALUES ($1, $2)`, notificationTable)

	_, err := n.executor(ctx).Exec(ctx, q, notification.ID, notification.DeviceToken)
	return err
}

func (n *NotificationRepository) Get(ctx context.Context, userID uuid.UUID) (domain.Notification, error) {
	q := fmt.Sprintf(`SELECT id, device_token FROM %s 
WHERE id = $1`, notificationTable)

	var notification domain.Notification
	row := n.executor(ctx).QueryRow(ctx, q, userID)
	err := row.Scan(&notification.ID, &notification.DeviceToken)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Notification{}, domain.ErrEntityIsNotExists
	}
	if err != nil {
		return domain.Notification{}, fmt.Errorf("can't get notification info: %w", err)
	}

	return notification, nil
}

func (n *NotificationRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	q := fmt.Sprintf(`DELETE FROM %s
WHERE id = $1`, notificationTable)

	_, err := n.executor(ctx).Exec(ctx, q, userID)
	if err != nil {
		return fmt.Errorf("can't delete notification in db: %w", err)
	}

	return nil
}
