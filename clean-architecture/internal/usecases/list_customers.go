package usecases

import "clean-architecture/internal/entities"

type ListCustomersInput struct {
	ActiveOnly bool
}

type CustomerListItem struct {
	CustomerID string
	Active     bool
}

type ListCustomersOutput struct {
	ActiveOnly bool
	Count      int
	Customers  []CustomerListItem
}

type ListCustomersInputBoundary interface {
	Execute(input ListCustomersInput) error
}

type ListCustomersOutputBoundary interface {
	Present(output ListCustomersOutput) error
}

type CustomerLister interface {
	List(activeOnly bool) ([]entities.Customer, error)
}

type ListCustomersInteractor struct {
	customers CustomerLister
	output    ListCustomersOutputBoundary
}

func NewListCustomersInteractor(customers CustomerLister, output ListCustomersOutputBoundary) ListCustomersInteractor {
	return ListCustomersInteractor{
		customers: customers,
		output:    output,
	}
}

func (uc ListCustomersInteractor) Execute(input ListCustomersInput) error {
	customers, err := uc.customers.List(input.ActiveOnly)
	if err != nil {
		return err
	}

	items := make([]CustomerListItem, 0, len(customers))
	for _, customer := range customers {
		items = append(items, CustomerListItem{
			CustomerID: customer.ID,
			Active:     customer.Active,
		})
	}

	return uc.output.Present(ListCustomersOutput{
		ActiveOnly: input.ActiveOnly,
		Count:      len(items),
		Customers:  items,
	})
}
