package memory

import "sync"

type IdempotencyStore struct {
	mu      sync.RWMutex
	entries map[string]string
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		entries: make(map[string]string),
	}
}

func (s *IdempotencyStore) Get(key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	status, ok := s.entries[key]
	return status, ok, nil
}

func (s *IdempotencyStore) Save(key string, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries[key] = status
	return nil
}
