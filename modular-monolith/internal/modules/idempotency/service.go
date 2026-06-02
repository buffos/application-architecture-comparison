package idempotency

type Service struct {
	store Store
}

func NewService(store Store) Service {
	return Service{
		store: store,
	}
}

func (s Service) Find(key string) (Result, bool, error) {
	return s.store.Find(key)
}

func (s Service) Save(key string, result Result) error {
	return s.store.Save(key, result)
}
