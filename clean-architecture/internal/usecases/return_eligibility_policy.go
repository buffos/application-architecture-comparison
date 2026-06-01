package usecases

import "clean-architecture/internal/entities"

type ReturnEligibilityPolicy interface {
	CanAccept(order entities.Order, request entities.ReturnRequest) (bool, error)
}
