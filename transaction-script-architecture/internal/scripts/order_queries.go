package scripts

import (
	"sort"

	"transaction-script-architecture/internal/data"
)

func GetOrder(store *data.Store, orderID string) (data.Order, error) {
	if store == nil {
		return data.Order{}, ErrStoreRequired
	}
	if orderID == "" {
		return data.Order{}, ErrOrderIDRequired
	}

	order, ok := store.Orders[orderID]
	if !ok {
		return data.Order{}, ErrOrderNotFound
	}

	return cloneOrder(order), nil
}

func ListOrders(store *data.Store, status string) ([]data.Order, error) {
	if store == nil {
		return nil, ErrStoreRequired
	}

	orders := make([]data.Order, 0, len(store.Orders))
	for _, order := range store.Orders {
		if status != "" && order.Status != status {
			continue
		}
		orders = append(orders, cloneOrder(order))
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].ID < orders[j].ID
	})

	return orders, nil
}

func cloneOrder(order data.Order) data.Order {
	order.Lines = append([]data.OrderLine(nil), order.Lines...)
	return order
}
