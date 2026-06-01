package memory

import "sync"

type IdempotencyStore struct {
	mu      sync.RWMutex
	records map[string]string
}

func NewIdempotencyStore() *IdempotencyStore {
	return &IdempotencyStore{
		records: make(map[string]string),
	}
}

func (s *IdempotencyStore) Find(commandName string, key string) (string, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resultID, ok := s.records[namespacedKey(commandName, key)]
	return resultID, ok, nil
}

func (s *IdempotencyStore) Save(commandName string, key string, resultID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.records[namespacedKey(commandName, key)] = resultID
	return nil
}

func namespacedKey(commandName string, key string) string {
	return commandName + ":" + key
}
