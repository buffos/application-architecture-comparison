package clock

import "time"

// Reader is the public time contract provided to components that need a
// timestamp without depending directly on the system clock.
type Reader interface {
	Now() time.Time
}

// Component provides the current system time at the composition boundary.
type Component struct{}

func NewComponent() *Component {
	return &Component{}
}

func (c *Component) Now() time.Time {
	return time.Now()
}

var _ Reader = (*Component)(nil)
