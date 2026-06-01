package usecases

import "errors"

var ErrIdempotencyKeyRequired = errors.New("idempotency key is required")

type IdempotencyStore interface {
	Find(commandName string, key string) (resultID string, found bool, err error)
	Save(commandName string, key string, resultID string) error
}
