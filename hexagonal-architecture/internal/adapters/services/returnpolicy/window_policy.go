package returnpolicy

import (
	"time"

	"hexagonal-architecture/internal/core/domain"
)

type WindowPolicy struct{}

func NewWindowPolicy() WindowPolicy {
	return WindowPolicy{}
}

func (WindowPolicy) CanAccept(request domain.ReturnRequest) (bool, error) {
	for _, line := range request.Lines {
		deadline := request.ShippedAt.AddDate(0, 0, line.ReturnWindowDays)
		if request.RequestedAt.After(endOfDay(deadline)) {
			return false, nil
		}
	}

	return true, nil
}

func endOfDay(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), value.Location())
}
