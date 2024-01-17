package service

import "sync"

type AttemptsRecorder struct {
	maxAttempts int64
	m           map[string]int64
	mu          sync.RWMutex
}

func newAttemptsRecorder(maxAttempts int64) *AttemptsRecorder {
	return &AttemptsRecorder{
		maxAttempts: maxAttempts,
		m:           make(map[string]int64, 10),
		mu:          sync.RWMutex{},
	}
}

func (i *AttemptsRecorder) addAttempt(email string) bool {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.m[email]++

	if i.m[email] > i.maxAttempts {
		delete(i.m, email)
		return false
	}

	return true
}

func (i *AttemptsRecorder) clearAttempts(email string) {
	i.mu.Lock()
	delete(i.m, email)
	i.mu.Unlock()
}
