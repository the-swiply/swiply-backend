package algorithm

import (
	"github.com/the-swiply/swiply-backend/randomcoffee/internal/domain"
)

type RandomCoffeeAlgorithm struct {
	cfg RandomCoffeeAlgorithmConfig
}

func NewRandomCoffeeAlgorithm(cfg RandomCoffeeAlgorithmConfig) *RandomCoffeeAlgorithm {
	return &RandomCoffeeAlgorithm{cfg: cfg}
}

func (r *RandomCoffeeAlgorithm) MatchUsers(meetings []domain.Meeting) []domain.Meeting {
	//TODO implement me
	panic("implement me")
}
