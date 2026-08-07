package customers

import (
	"errors"
	"sort"

	"rich-domain-model-architecture/internal/domain/customer"
)

var ErrCustomerNotFound = errors.New("customer not found")

type Reader interface {
	GetCustomer(id string) (CustomerDetails, error)
	ListCustomers(active *bool) []CustomerSummary
}

type CustomerDetails struct {
	CustomerID  string
	Tier        string
	PaymentTerm string
	Active      bool
}

type CustomerSummary struct {
	CustomerID string
	Tier       string
	Active     bool
}

type InMemoryReader struct {
	customers map[string]CustomerDetails
}

func NewInMemoryReader() *InMemoryReader {
	return &InMemoryReader{customers: make(map[string]CustomerDetails)}
}

func (reader *InMemoryReader) Save(value customer.Customer) {
	reader.customers[string(value.ID())] = CustomerDetails{
		CustomerID:  string(value.ID()),
		Tier:        string(value.Tier()),
		PaymentTerm: string(value.PaymentTerms()),
		Active:      value.Active(),
	}
}

func (reader *InMemoryReader) GetCustomer(id string) (CustomerDetails, error) {
	value, ok := reader.customers[id]
	if !ok {
		return CustomerDetails{}, ErrCustomerNotFound
	}
	return value, nil
}

func (reader *InMemoryReader) ListCustomers(active *bool) []CustomerSummary {
	result := make([]CustomerSummary, 0, len(reader.customers))
	for _, value := range reader.customers {
		if active != nil && value.Active != *active {
			continue
		}
		result = append(result, CustomerSummary{CustomerID: value.CustomerID, Tier: value.Tier, Active: value.Active})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].CustomerID < result[j].CustomerID })
	return result
}

var _ Reader = (*InMemoryReader)(nil)
