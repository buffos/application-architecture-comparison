package ports

import "hexagonal-architecture/internal/core/domain"

type ReturnEligibilityPolicy interface {
	CanAccept(request domain.ReturnRequest) (bool, error)
}
