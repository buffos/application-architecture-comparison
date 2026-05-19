package timeadapter

import "time"

type FixedClock struct {
	current time.Time
}

func NewFixedClock(current time.Time) FixedClock {
	return FixedClock{current: current}
}

func (c FixedClock) Now() time.Time {
	return c.current
}
