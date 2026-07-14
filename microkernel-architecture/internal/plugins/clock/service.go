package clock

import "time"

type Service struct{}

func NewService() Service {
	return Service{}
}

func (s Service) Now() time.Time {
	return time.Now().UTC()
}
