package ports

type IdempotencyStore interface {
	Seen(scope, key string) (bool, error)
	Remember(scope, key, resourceID string) error
	ResourceID(scope, key string) (string, error)
}
