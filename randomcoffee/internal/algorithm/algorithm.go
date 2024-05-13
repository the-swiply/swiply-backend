package algorithm

import (
	"sync"
	"time"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/domain"
)

type RandomCoffeeAlgorithm struct {
	cfg RandomCoffeeAlgorithmConfig
}

func NewRandomCoffeeAlgorithm(cfg RandomCoffeeAlgorithmConfig) *RandomCoffeeAlgorithm {
	return &RandomCoffeeAlgorithm{cfg: cfg}
}

func (r *RandomCoffeeAlgorithm) MatchUsers(meetings []domain.Meeting) []domain.Meeting {
	groupedByOrganizations := groupBy(meetings, func(item domain.Meeting) int64 {
		return item.OrganizationID
	})

	var (
		wg sync.WaitGroup

		mu           sync.Mutex
		matchedUsers []domain.Meeting
	)

	for _, meetingsByOrganization := range groupedByOrganizations {
		wg.Add(1)
		go func(meetings []domain.Meeting) {
			defer wg.Done()

			answer := r.matchUsersInOrganization(meetings)

			mu.Lock()
			matchedUsers = append(matchedUsers, answer...)
			mu.Unlock()
		}(meetingsByOrganization)
	}

	wg.Wait()

	return matchedUsers
}

func (r *RandomCoffeeAlgorithm) matchUsersInOrganization(meetings []domain.Meeting) []domain.Meeting {
	var reservedUsers []domain.Meeting
	reserved := int(24 * time.Hour / r.cfg.Interval)

	groupedByAvgMeetingTime := groupBy(meetings, func(item domain.Meeting) int64 {
		return item.Start.Add((item.End.Sub(item.Start) / 2) / r.cfg.Interval * r.cfg.Interval).Unix()
	})

	for _, meetingByAvgMeetingTime := range groupedByAvgMeetingTime {
		sortByMeetingTime := make([][]domain.Meeting, reserved, reserved)
		for _, meeting := range meetingByAvgMeetingTime {
			sortByMeetingTime[int(meeting.End.Sub(meeting.Start)/r.cfg.Interval)-1] =
				append(sortByMeetingTime[int(meeting.End.Sub(meeting.Start)/r.cfg.Interval)-1], meeting)
		}

		var resUsers []domain.Meeting
		for i := reserved - 1; i >= 0; i-- {
			for j := len(sortByMeetingTime[i]) - 1; 0 <= j && len(resUsers) < reserved; j-- {
				resUsers = append(resUsers, sortByMeetingTime[i][j])
				sortByMeetingTime[i] = sortByMeetingTime[i][:len(sortByMeetingTime[i])-1]
			}
		}

		reservedUsers = append(reservedUsers, resUsers...)

		for i := 0; i < reserved; i++ {

		}
	}
}

func groupBy[T any, E comparable](items []T, key func(item T) E) (answer [][]T) {
	mp := make(map[E][]T)
	for _, item := range items {
		mp[key(item)] = append(mp[key(item)], item)
	}

	for _, value := range mp {
		answer = append(answer, value)
	}

	return
}
