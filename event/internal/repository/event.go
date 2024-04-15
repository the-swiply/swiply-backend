package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/event/internal/domain"
)

const (
	eventTable = "event"
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
	q := fmt.Sprintf(`INSERT INTO %s (owner, title, description, date)
VALUES ($1, $2, $3, $4) RETURNING id`, eventTable)

	var id int64
	row := e.db.QueryRow(ctx, q, event.Owner, event.Title, event.Description, event.Date)

	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("can't insert event to db: %w", err)
	}

	return id, err
}

func (e *EventRepository) UpdateEvent(ctx context.Context, event domain.Event) error {
	q := fmt.Sprintf(`UPDATE %s
SET title = $1,
    description = $2,
    date = $3
WHERE id = $4 AND owner = $5`, eventTable)

	_, err := e.db.Exec(ctx, q, event.Title, event.Description, event.Date, event.ID, event.Owner)
	if err != nil {
		return fmt.Errorf("can't update event in db: %w", err)
	}

	return nil
}

func (e *EventRepository) GetEvents(ctx context.Context, owner uuid.UUID) ([]domain.Event, error) {
	q := fmt.Sprintf(`SELECT id, owner, members, title, description, date FROM %s 
WHERE owner = $1
ORDER BY id DESC`, eventTable)

	rows, err := e.db.Query(ctx, q, owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Event])
}
