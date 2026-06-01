package memory

import (
	"sort"
	"sync"

	"clean-architecture/internal/entities"
)

type OrderGateway struct {
	mu     sync.RWMutex
	orders map[string]entities.Order
}

func NewOrderGateway() *OrderGateway {
	return &OrderGateway{
		orders: make(map[string]entities.Order),
	}
}

func (g *OrderGateway) Save(order entities.Order) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.orders[order.ID] = order
	return nil
}

func (g *OrderGateway) FindByID(id string) (entities.Order, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	order, ok := g.orders[id]
	if !ok {
		return entities.Order{}, entities.ErrQuoteNotFound
	}

	return order, nil
}

func (g *OrderGateway) ListByStatus(status string) ([]entities.Order, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	orders := make([]entities.Order, 0, len(g.orders))
	for _, order := range g.orders {
		if status != "" && order.Status != status {
			continue
		}

		orders = append(orders, order)
	}

	sort.Slice(orders, func(i int, j int) bool {
		return orders[i].ID < orders[j].ID
	})

	return orders, nil
}
