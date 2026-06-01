package returneligibility

import "onion-architecture/internal/domain"

type ReasonPolicy struct{}

func NewReasonPolicy() ReasonPolicy {
	return ReasonPolicy{}
}

func (p ReasonPolicy) IsEligible(request domain.ReturnRequest, order domain.Order) (bool, error) {
	if request.Reason == "outside return window" {
		return false, nil
	}

	return true, nil
}
