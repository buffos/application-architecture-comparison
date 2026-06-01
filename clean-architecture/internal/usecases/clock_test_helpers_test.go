package usecases

import "time"

type stubClock struct {
	now time.Time
}

func (c stubClock) Now() time.Time {
	return c.now
}
