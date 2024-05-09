package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/domain"
)

type MeetingRepository interface {
	Create(ctx context.Context, meeting domain.Meeting) error
	Get(ctx context.Context, meetingID uuid.UUID) (domain.Meeting, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]domain.Meeting, error)
	Update(ctx context.Context, meeting domain.Meeting) error
	Delete(ctx context.Context, id, ownerID uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.MeetingStatus) error
	ListRoundMeetings(ctx context.Context, start time.Time) ([]domain.Meeting, error)
	UpdateMember(ctx context.Context, id, memberID uuid.UUID) error
}

type MeetingService struct {
	config *MeetingConfig
	repo   MeetingRepository
}

func NewMeetingService(config *MeetingConfig, repo MeetingRepository) *MeetingService {
	return &MeetingService{config: config, repo: repo}
}

func (m *MeetingService) Create(ctx context.Context, meeting domain.Meeting) error {
	return m.repo.Create(ctx, meeting)
}

func (m *MeetingService) Get(ctx context.Context, meetingID uuid.UUID) (domain.Meeting, error) {
	return m.repo.Get(ctx, meetingID)
}

func (m *MeetingService) List(ctx context.Context, ownerID uuid.UUID) ([]domain.Meeting, error) {
	return m.repo.List(ctx, ownerID)
}

func (m *MeetingService) Update(ctx context.Context, meeting domain.Meeting) error {
	return m.repo.Update(ctx, meeting)
}

func (m *MeetingService) Delete(ctx context.Context, id, ownerID uuid.UUID) error {
	return m.repo.Delete(ctx, id, ownerID)
}
