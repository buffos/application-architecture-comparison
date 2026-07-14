package memory

import (
	"sync"

	"microkernel-architecture/internal/kernel"
)

type IdempotencyStore struct {
	mu      sync.RWMutex
	results map[string]kernel.IdempotencyResult
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		results: make(map[string]kernel.IdempotencyResult),
	}
}

func (s *IdempotencyStore) Find(key string) (kernel.IdempotencyResult, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result, ok := s.results[key]
	return result, ok, nil
}

func (s *IdempotencyStore) Save(key string, result kernel.IdempotencyResult) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.results[key] = result
	return nil
}
