package orders

type GetOrderQuery struct {
	OrderID string
}

type OrderDetails struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type ListOrdersQuery struct {
	Status string
}

func (s Service) GetOrder(query GetOrderQuery) (OrderDetails, error) {
	order, err := s.orders.FindByID(query.OrderID)
	if err != nil {
		return OrderDetails{}, err
	}

	return OrderDetails{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}

func (s Service) ListOrders(query ListOrdersQuery) ([]OrderDetails, error) {
	orders, err := s.orders.ListByStatus(query.Status)
	if err != nil {
		return nil, err
	}

	list := make([]OrderDetails, 0, len(orders))
	for _, order := range orders {
		list = append(list, OrderDetails{
			OrderID:    order.ID,
			QuoteID:    order.QuoteID,
			CustomerID: order.CustomerID,
			Status:     order.Status,
			LineCount:  len(order.Lines),
		})
	}

	return list, nil
}
