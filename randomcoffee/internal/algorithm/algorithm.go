package algorithm

import (
	"math/rand"
	"sync"
	"time"

	"github.com/the-swiply/swiply-backend/randomcoffee/internal/domain"
	"github.com/the-swiply/swiply-backend/randomcoffee/pkg/edmonds"
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
	var (
		reservedUsers []domain.Meeting
		matchedUsers  []domain.Meeting
	)
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

		weights := make([]*weight, reserved, reserved)
		weights[reserved-1].weight = 1
		for i := reserved - 2; i >= 0; i-- {
			weights[i].weight = 2 * weights[i+1].weight
		}
		for i, value := range sortByMeetingTime {
			weights[i].members = len(value)
		}

		for i := 0; i < reserved; i++ {
			for 0 < len(sortByMeetingTime[i]) {
				member1 := sortByMeetingTime[i][len(sortByMeetingTime[i])-1]
				sortByMeetingTime[i] = sortByMeetingTime[i][:len(sortByMeetingTime[i])-1]
				weights[i].members--
				if weights[i].members == 0 {
					weights[i].weight = 0
				}

				group := randomGroup(weights)
				if group == -1 {
					reservedUsers = append(reservedUsers, member1)
					break
				}

				member2 := sortByMeetingTime[group][len(sortByMeetingTime[group])-1]
				sortByMeetingTime[group] = sortByMeetingTime[group][:len(sortByMeetingTime[group])-1]
				weights[group].members--
				if weights[group].members == 0 {
					weights[group].weight = 0
				}

				member1.MemberID = member2.OwnerID
				member2.MemberID = member1.OwnerID
				matchedUsers = append(matchedUsers, member1, member2)
			}
		}
	}

	graph := edmonds.NewGraph(func(meeting domain.Meeting) string {
		return meeting.OwnerID.String()
	})

	for _, userA := range reservedUsers {
		for _, userB := range reservedUsers {
			if !userB.Start.Add(r.cfg.Interval).After(userA.End) {
				graph.AddEdge(userA, userB)
			}
		}
	}

	for _, pair := range graph.MatchPairs() {
		pair.First.MemberID = pair.Second.OwnerID
		pair.Second.MemberID = pair.First.OwnerID
		matchedUsers = append(matchedUsers, pair.First, pair.Second)
	}

	return matchedUsers
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

type weight struct {
	members int
	weight  int64
}

func randomGroup(weights []*weight) int {
	var (
		sum   int64
		total int
	)

	for _, value := range weights {
		total += value.members
		sum += value.weight
	}

	if total == 0 {
		return -1
	}

	randSum := rand.Int63n(sum) + 1
	for i, value := range weights {
		if randSum <= value.weight && value.weight != 0 {
			return i
		} else {
			randSum -= value.weight
		}
	}

	return len(weights) - 1
}
