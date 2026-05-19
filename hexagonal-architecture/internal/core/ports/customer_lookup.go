package ports

import "hexagonal-architecture/internal/core/domain"

type CustomerLookup interface {
	FindByID(id string) (domain.Customer, error)
	List(activeOnly bool) ([]domain.Customer, error)
}
