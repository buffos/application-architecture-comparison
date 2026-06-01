package application

type IdempotencyStore interface {
	Get(key string) (string, bool, error)
	Save(key string, status string) error
}
