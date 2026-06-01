package returneligibility

import "onion-architecture/internal/domain"

type WindowPolicy struct{}

func NewWindowPolicy() WindowPolicy {
	return WindowPolicy{}
}

func (p WindowPolicy) IsEligible(request domain.ReturnRequest, order domain.Order) (bool, error) {
	if order.ShippedAt.IsZero() {
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
