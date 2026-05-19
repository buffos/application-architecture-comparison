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

func (s *IdempotencyStore) Seen(scope, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.entries[scope+"::"+key]
	return ok, nil
}

func (s *IdempotencyStore) Remember(scope, key, resourceID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entries[scope+"::"+key] = resourceID
	return nil
}

func (s *IdempotencyStore) ResourceID(scope, key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.entries[scope+"::"+key], nil
}
