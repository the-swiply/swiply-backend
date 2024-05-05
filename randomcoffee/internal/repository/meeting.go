package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/domain"
)

const meetingTable = "meeting"

type MeetingRepository struct {
	db *pgxpool.Pool
}

func NewMeetingRepository(db *pgxpool.Pool) *MeetingRepository {
	return &MeetingRepository{
		db: db,
	}
}

func (m *MeetingRepository) Create(ctx context.Context, meeting domain.Meeting) error {
	q := fmt.Sprintf()
}
