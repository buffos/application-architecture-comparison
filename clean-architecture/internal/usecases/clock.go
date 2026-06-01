package usecases

import "time"

type Clock interface {
	Now() time.Time
}
