package idempotency

import "microkernel-architecture/internal/kernel"

type Store interface {
	Find(key string) (kernel.IdempotencyResult, bool, error)
	Save(key string, result kernel.IdempotencyResult) error
}

type Service struct {
	store Store
}

func NewService(store Store) Service {
	return Service{
		store: store,
	}
}

func (s Service) Find(key string) (kernel.IdempotencyResult, bool, error) {
	return s.store.Find(key)
}

func (s Service) Save(key string, result kernel.IdempotencyResult) error {
	return s.store.Save(key, result)
}
