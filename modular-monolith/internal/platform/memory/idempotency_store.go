package memory

import (
	"sync"

	"modular-monolith/internal/modules/idempotency"
)

type IdempotencyStore struct {
	mu      sync.RWMutex
	results map[string]idempotency.Result
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		results: make(map[string]idempotency.Result),
	}
}

func (s *IdempotencyStore) Find(key string) (idempotency.Result, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, ok := s.results[key]
	return result, ok, nil
}

func (s *IdempotencyStore) Save(key string, result idempotency.Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.results[key] = result
	return nil
}
