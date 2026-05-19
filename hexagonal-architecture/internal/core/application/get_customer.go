package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type GetCustomerUseCase struct {
	customers ports.CustomerLookup
}

func NewGetCustomerUseCase(customers ports.CustomerLookup) GetCustomerUseCase {
	return GetCustomerUseCase{customers: customers}
}

func (uc GetCustomerUseCase) Execute(id string) (domain.Customer, error) {
	return uc.customers.FindByID(id)
}
