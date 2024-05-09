package service

import (
	"context"
	"time"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/domain"
)

type RandomCoffeeAlgorithm interface {
	MatchUsers(meetings []domain.Meeting) []domain.Meeting
}

type RandomCoffeeService struct {
	config    *RandomCoffeeConfig
	algorithm RandomCoffeeAlgorithm
	repo      MeetingRepository
}

func NewRandomCoffeeService(config *RandomCoffeeConfig, algorithm RandomCoffeeAlgorithm, repo MeetingRepository) *RandomCoffeeService {
	return &RandomCoffeeService{config: config, algorithm: algorithm, repo: repo}
}

func (r *RandomCoffeeService) Schedule(ctx context.Context) error {
	meetings, err := r.repo.ListRoundMeetings(ctx, time.Now().Add(24*time.Hour))
	if err != nil {
		return err
	}

	for _, meeting := range meetings {
		_ = r.repo.UpdateStatus(ctx, meeting.ID, domain.MeetingStatusScheduling)
	}

	meetings = r.algorithm.MatchUsers(meetings)

	for _, meeting := range meetings {
		_ = r.repo.UpdateStatus(ctx, meeting.ID, domain.MeetingStatusScheduled)
		_ = r.repo.UpdateMember(ctx, meeting.ID, meeting.MemberID)
	}

	return nil
}
