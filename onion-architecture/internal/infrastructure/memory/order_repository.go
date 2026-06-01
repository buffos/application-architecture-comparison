package memory

import (
	"sync"

	"onion-architecture/internal/domain"
)

type OrderRepository struct {
	mu     sync.RWMutex
	orders map[string]domain.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		orders: make(map[string]domain.Order),
	}
}

func (r *OrderRepository) Save(order domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[order.ID] = order
	return nil
}

func (r *OrderRepository) FindByID(id string) (domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[id]
	if !ok {
		return domain.Order{}, domain.ErrOrderNotFound
	}

	return order, nil
}

func (r *OrderRepository) ListByStatus(status string) ([]domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]domain.Order, 0)
	for _, order := range r.orders {
		if status == "" {
			result = append(result, order)
			continue
		}

		if order.Status == status {
			result = append(result, order)
		}
	}

	return result, nil
}
