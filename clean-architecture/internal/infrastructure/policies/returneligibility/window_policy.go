package returneligibility

import (
	"clean-architecture/internal/entities"
)

type WindowPolicy struct{}

func NewWindowPolicy() WindowPolicy {
	return WindowPolicy{}
}

func (p WindowPolicy) CanAccept(order entities.Order, request entities.ReturnRequest) (bool, error) {
	if order.ShippedAt == nil {
		return false, nil
	}

	for _, line := range order.Lines {
		if line.ReturnWindowDays <= 0 {
			return false, nil
		}

		deadline := order.ShippedAt.AddDate(0, 0, line.ReturnWindowDays)
		if request.RequestedAt.After(deadline) {
			return false, nil
		}
	}

	return true, nil
}
