package time

import stdtime "time"

type SystemClock struct{}

func NewSystemClock() SystemClock {
	return SystemClock{}
}

func (c SystemClock) Now() stdtime.Time {
	return stdtime.Now()
}
