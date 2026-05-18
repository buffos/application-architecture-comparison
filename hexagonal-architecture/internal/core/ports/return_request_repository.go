package ports

import "hexagonal-architecture/internal/core/domain"

type ReturnRequestRepository interface {
	Save(request domain.ReturnRequest) error
	FindByID(id string) (domain.ReturnRequest, error)
}
