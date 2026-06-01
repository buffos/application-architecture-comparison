package usecases

import "clean-architecture/internal/entities"

type GetCustomerInput struct {
	CustomerID string
}

type GetCustomerOutput struct {
	CustomerID string
	Active     bool
}

type GetCustomerInputBoundary interface {
	Execute(input GetCustomerInput) error
}

type GetCustomerOutputBoundary interface {
	Present(output GetCustomerOutput) error
}

type CustomerReader interface {
	FindByID(id string) (entities.Customer, error)
}

type GetCustomerInteractor struct {
	customers CustomerReader
	output    GetCustomerOutputBoundary
}

func NewGetCustomerInteractor(customers CustomerReader, output GetCustomerOutputBoundary) GetCustomerInteractor {
	return GetCustomerInteractor{
		customers: customers,
		output:    output,
	}
}

func (uc GetCustomerInteractor) Execute(input GetCustomerInput) error {
	customer, err := uc.customers.FindByID(input.CustomerID)
	if err != nil {
		return err
	}

	return uc.output.Present(GetCustomerOutput{
		CustomerID: customer.ID,
		Active:     customer.Active,
	})
}
