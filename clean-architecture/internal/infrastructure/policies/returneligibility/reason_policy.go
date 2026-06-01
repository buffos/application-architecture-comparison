package returneligibility

import (
	"strings"

	"clean-architecture/internal/entities"
)

type ReasonPolicy struct{}

func NewReasonPolicy() ReasonPolicy {
	return ReasonPolicy{}
}

func (p ReasonPolicy) CanAccept(order entities.Order, request entities.ReturnRequest) (bool, error) {
	if strings.EqualFold(strings.TrimSpace(request.Reason), "outside return window") {
		return false, nil
	}

	return true, nil
}
