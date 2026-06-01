package memory

import (
	"sync"

	"clean-architecture/internal/entities"
)

type CustomerGateway struct {
	mu        sync.RWMutex
	customers map[string]entities.Customer
}

func NewCustomerGateway() *CustomerGateway {
	return &CustomerGateway{
		customers: make(map[string]entities.Customer),
	}
}

func (g *CustomerGateway) Save(customer entities.Customer) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.customers[customer.ID] = customer
	return nil
}

func (g *CustomerGateway) FindByID(id string) (entities.Customer, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	customer, ok := g.customers[id]
	if !ok {
		return entities.Customer{}, entities.ErrCustomerNotFound
	}

	return customer, nil
}
