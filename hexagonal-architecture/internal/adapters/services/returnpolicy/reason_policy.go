package returnpolicy

import (
	"strings"

	"hexagonal-architecture/internal/core/domain"
)

type ReasonPolicy struct{}

func NewReasonPolicy() ReasonPolicy {
	return ReasonPolicy{}
}

func (ReasonPolicy) CanAccept(request domain.ReturnRequest) (bool, error) {
	if strings.Contains(strings.ToLower(request.Reason), "outside return window") {
		return false, nil
	}

	return true, nil
}
