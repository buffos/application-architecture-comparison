package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ListCustomersUseCase struct {
	customers ports.CustomerLookup
}

func NewListCustomersUseCase(customers ports.CustomerLookup) ListCustomersUseCase {
	return ListCustomersUseCase{customers: customers}
}

func (uc ListCustomersUseCase) Execute(activeOnly bool) ([]domain.Customer, error) {
	return uc.customers.List(activeOnly)
}
