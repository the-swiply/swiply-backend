package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/chat/internal/domain"
)

const (
	chatTable    = "chat"
	messageTable = "message"
)

type ChatRepository struct {
	db *pgxpool.Pool
}

func NewChatRepository(db *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (c *ChatRepository) GetChatMembers(ctx context.Context, chatID int64) ([]uuid.UUID, error) {
	q := fmt.Sprintf(`SELECT members FROM %s WHERE id = $1`, chatTable)

	row := c.db.QueryRow(ctx, q, chatID)

	var members []uuid.UUID
	err := row.Scan(&members)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrEntityIsNotExists
	}
	if err != nil {
		return nil, fmt.Errorf("can't get chat members: %w", err)
	}

	return members, nil
}

func (c *ChatRepository) SaveMessage(ctx context.Context, msg domain.ChatMessage) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, chat_id, id_in_chat, send_time, content) VALUES ($1, $2, $3, $4, $5)`,
		messageTable)

	_, err := c.db.Exec(ctx, q, msg.ID, msg.ChatID, msg.IDInChat, msg.SendTime, msg.Content)
	return err
}
