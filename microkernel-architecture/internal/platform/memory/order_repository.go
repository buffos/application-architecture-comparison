package memory

import (
	"sync"

	"microkernel-architecture/internal/plugins/orders"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]orders.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]orders.Order),
	}
}

func (r *OrderRepository) Save(order orders.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepository) FindByID(id string) (orders.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return orders.Order{}, orders.ErrOrderNotFound
	}

	return order, nil
}
