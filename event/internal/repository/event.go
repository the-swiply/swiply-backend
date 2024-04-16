package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/event/internal/domain"
)

const (
	eventTable           = "event"
	eventUserStatusTable = "event_user_status"

	statusJoinRequest = "join_request"
	statusMember      = "member"
)

type EventRepository struct {
	db *pgxpool.Pool
}

func NewEventRepository(db *pgxpool.Pool) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (e *EventRepository) CreateEvent(ctx context.Context, event domain.Event) (int64, error) {
	tx, err := e.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return 0, fmt.Errorf("can't begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	q1 := fmt.Sprintf(`INSERT INTO %s (owner, title, description, date)
VALUES ($1, $2, $3, $4) RETURNING id`, eventTable)

	var id int64
	row := e.db.QueryRow(ctx, q1, event.Owner, event.Title, event.Description, event.Date)

	err = row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't insert event to db: %w", err)
	}

	q2 := fmt.Sprintf(`INSERT INTO %s (user_id, event_id, status) VALUES ($1, $2, $3)`, eventUserStatusTable)
	_, err = e.db.Exec(ctx, q2, event.Owner, id, statusMember)
	if err != nil {
		return 0, fmt.Errorf("can't insert status: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("can't commit tx: %w", err)
	}

	return id, err
}

func (e *EventRepository) UpdateEvent(ctx context.Context, event domain.Event) error {
	q := fmt.Sprintf(`UPDATE %s
SET title = $1,
    description = $2,
    chat_id = $3,
    date = $4
WHERE id = $5 AND owner = $6`, eventTable)

	_, err := e.db.Exec(ctx, q, event.Title, event.Description, event.ChatID, event.Date, event.ID, event.Owner)
	if err != nil {
		return fmt.Errorf("can't update event in db: %w", err)
	}

	return nil
}

func (e *EventRepository) GetEvents(ctx context.Context, owner uuid.UUID) ([]domain.Event, error) {
	q := fmt.Sprintf(`SELECT id, owner, title, description, chat_id, date FROM %s 
WHERE owner = $1
ORDER BY date DESC`, eventTable)

	rows, err := e.db.Query(ctx, q, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Event])
}

func (e *EventRepository) GetEventByID(ctx context.Context, id int64) (domain.Event, error) {
	q := fmt.Sprintf(`SELECT id, owner, title, description, chat_id, date FROM %s 
WHERE id = $1`, eventTable)

	var ev domain.Event
	row := e.db.QueryRow(ctx, q, id)
	err := row.Scan(&ev.ID, &ev.Owner, &ev.Title, &ev.Description, &ev.ChatID, &ev.Date)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Event{}, domain.ErrEntityIsNotExists
	}
	if err != nil {
		return domain.Event{}, fmt.Errorf("can't get event: %w", err)
	}

	return ev, nil
}

func (e *EventRepository) JoinEvent(ctx context.Context, eventID int64, userID uuid.UUID) error {
	q := fmt.Sprintf(`INSERT INTO %s (user_id, event_id, status) VALUES ($1, $2, $3) ON CONFLICT (user_id, event_id) DO NOTHING`,
		eventUserStatusTable)

	_, err := e.db.Exec(ctx, q, userID, eventID, statusJoinRequest)
	if err != nil {
		return fmt.Errorf("can't update event in db: %w", err)
	}

	return nil
}

func (e *EventRepository) AcceptEventJoin(ctx context.Context, eventID int64, owner, userID uuid.UUID) error {
	q := fmt.Sprintf(`UPDATE %s SET status = $1
          WHERE event_id = $3 AND user_id = $2 AND status = $5 AND (SELECT owner FROM %s WHERE id = $3) = $4`,
		eventUserStatusTable, eventTable)

	_, err := e.db.Exec(ctx, q, statusMember, userID, eventID, owner, statusJoinRequest)
	if err != nil {
		return fmt.Errorf("can't update event in db: %w", err)
	}

	return nil
}
