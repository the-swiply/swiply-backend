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
	q := fmt.Sprintf(`INSERT INTO %s (id, chat_id, id_in_chat, "from", send_time, content)
VALUES ($1, $2, $3, $4, $5, $6)`, messageTable)

	_, err := c.db.Exec(ctx, q, msg.ID, msg.ChatID, msg.IDInChat, msg.From, msg.SendTime, msg.Content)
	return err
}

func (c *ChatRepository) GetNextMessages(ctx context.Context, chatID int64, start int64, limit int64) ([]domain.ChatMessage, error) {
	q := fmt.Sprintf(`SELECT id, chat_id, id_in_chat, "from", send_time, content FROM %s 
WHERE chat_id = $1 AND id_in_chat > $2
ORDER BY id_in_chat
LIMIT $3`, messageTable)

	rows, err := c.db.Query(ctx, q, chatID, start, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.ChatMessage])
}

func (c *ChatRepository) GetPreviousMessages(ctx context.Context, chatID int64, start int64, limit int64) ([]domain.ChatMessage, error) {
	q := fmt.Sprintf(`SELECT id, chat_id, id_in_chat, "from", send_time, content FROM %s 
WHERE chat_id = $1 AND id_in_chat < $2
ORDER BY id_in_chat
LIMIT $3`, messageTable)

	rows, err := c.db.Query(ctx, q, chatID, start, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.ChatMessage])
}

func (c *ChatRepository) GetUserChats(ctx context.Context, userID uuid.UUID) ([]domain.Chat, error) {
	q := fmt.Sprintf(`SELECT id, members FROM %s
WHERE $1 = ANY(members)`, chatTable)

	rows, err := c.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Chat])
}

func (c *ChatRepository) RemoveUserFromChatMembers(ctx context.Context, userID uuid.UUID, chatID int64) error {
	q := fmt.Sprintf(`UPDATE %s
SET members = array_remove(members, $1)
WHERE id = $2`, chatTable)

	_, err := c.db.Exec(ctx, q, userID, chatID)
	return err
}

func (c *ChatRepository) CreateChat(ctx context.Context, members []uuid.UUID) error {
	q := fmt.Sprintf(`INSERT INTO %s (members)
VALUES ($1)`, chatTable)

	_, err := c.db.Exec(ctx, q, members)

	return err
}
