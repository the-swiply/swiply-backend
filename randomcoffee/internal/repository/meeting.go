package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/the-swiply/swiply-backend/pkg/houston/dobby"

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

func (m *MeetingRepository) executor(ctx context.Context) dobby.Executor {
	tx := dobby.ExtractPGXTx(ctx)
	if tx != nil {
		return tx
	}

	return m.db
}

func (m *MeetingRepository) Create(ctx context.Context, meeting domain.Meeting) error {
	q := fmt.Sprintf(`INSERT INTO %s (id, owner_id, member_id, "start", "end", organization_id, status, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, meetingTable)

	_, err := m.executor(ctx).Exec(ctx, q, meeting.ID, meeting.OwnerID, meeting.MemberID, meeting.Start, meeting.End,
		meeting.OrganizationID, meeting.Status, meeting.CreatedAt)
	return err
}

func (m *MeetingRepository) Get(ctx context.Context, meetingID uuid.UUID) (domain.Meeting, error) {
	q := fmt.Sprintf(`SELECT id, owner_id, member_id, "start", "end", organization_id, status, created_at FROM %s 
WHERE id = $1`, meetingTable)

	var meeting domain.Meeting
	row := m.executor(ctx).QueryRow(ctx, q, meetingID)
	err := row.Scan(&meeting.ID, &meeting.OwnerID, &meeting.MemberID, &meeting.Start, &meeting.End,
		&meeting.OrganizationID, &meeting.Status, &meeting.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Meeting{}, domain.ErrEntityIsNotExists
	}
	if err != nil {
		return domain.Meeting{}, fmt.Errorf("can't get meeting: %w", err)
	}

	return meeting, nil
}

func (m *MeetingRepository) List(ctx context.Context, ownerID uuid.UUID) ([]domain.Meeting, error) {
	q := fmt.Sprintf(`SELECT id, owner_id, member_id, "start", "end", organization_id, status, created_at FROM %s 
WHERE owner_id = $1`, meetingTable)

	rows, err := m.executor(ctx).Query(ctx, q, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Meeting])
}

func (m *MeetingRepository) Update(ctx context.Context, meeting domain.Meeting) error {
	q := fmt.Sprintf(`UPDATE %s
SET "start" = $1,
    "end" = $2,
    organization_id = $3,
WHERE id = $4 AND owner_id = $5 AND status = $6`, meetingTable)

	_, err := m.executor(ctx).Exec(ctx, q, meeting.Start, meeting.End, meeting.OrganizationID,
		meeting.ID, meeting.OwnerID, domain.MeetingStatusAwaitingSchedule)
	if err != nil {
		return fmt.Errorf("can't update meeting in db: %w", err)
	}

	return nil
}

func (m *MeetingRepository) Delete(ctx context.Context, id, ownerID uuid.UUID) error {
	q := fmt.Sprintf(`DELETE FROM %s
WHERE id = $1 AND owner_id = $2 AND status = $3`, meetingTable)

	_, err := m.executor(ctx).Exec(ctx, q, id, ownerID, domain.MeetingStatusAwaitingSchedule)
	if err != nil {
		return fmt.Errorf("can't delete meeting in db: %w", err)
	}

	return nil
}

func (m *MeetingRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.MeetingStatus) error {
	q := fmt.Sprintf(`UPDATE %s
SET status = $1
WHERE id = $2`, meetingTable)

	_, err := m.executor(ctx).Exec(ctx, q, status, id)
	if err != nil {
		return fmt.Errorf("can't update meeting status in db: %w", err)
	}

	return nil
}

func (m *MeetingRepository) ListRoundMeetings(ctx context.Context, start time.Time) ([]domain.Meeting, error) {
	q := fmt.Sprintf(`SELECT id, owner_id, member_id, "start", "end", organization_id, status, created_at FROM %s 
WHERE "start" >= $1 AND "end" < $2 AND status = $3`, meetingTable)

	rows, err := m.executor(ctx).Query(ctx, q, start, start.Add(24*time.Hour), domain.MeetingStatusAwaitingSchedule)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[domain.Meeting])
}

func (m *MeetingRepository) UpdateMember(ctx context.Context, id, memberID uuid.UUID) error {
	q := fmt.Sprintf(`UPDATE %s
SET member_id = $1
WHERE id = $2`, meetingTable)

	_, err := m.executor(ctx).Exec(ctx, q, memberID, id)
	if err != nil {
		return fmt.Errorf("can't update meeting member in db: %w", err)
	}

	return nil
}
