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

	for _, requestLine := range request.Lines {
		matched := false
		for _, orderLine := range order.Lines {
			if orderLine.SKU != requestLine.SKU {
				continue
			}

			if orderLine.ReturnWindowDays <= 0 {
				return false, nil
			}

			deadline := order.ShippedAt.AddDate(0, 0, orderLine.ReturnWindowDays)
			if request.RequestedAt.After(deadline) {
				return false, nil
			}

			matched = true
			break
		}

		if !matched {
			return false, nil
		}
	}

	return true, nil
}
