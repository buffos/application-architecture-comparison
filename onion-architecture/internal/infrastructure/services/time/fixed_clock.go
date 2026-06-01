package time

import stdtime "time"

type FixedClock struct {
	Current stdtime.Time
}

func NewFixedClock(current stdtime.Time) FixedClock {
	return FixedClock{Current: current}
}

func (c FixedClock) Now() stdtime.Time {
	return c.Current
}
