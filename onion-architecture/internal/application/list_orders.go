package application

type ListOrdersQuery struct {
	Status string
}

type ListOrdersService struct {
	orders OrderFinder
}

func NewListOrdersService(orders OrderFinder) ListOrdersService {
	return ListOrdersService{
		orders: orders,
	}
}

func (s ListOrdersService) Execute(query ListOrdersQuery) ([]OrderDetails, error) {
	orders, err := s.orders.ListByStatus(query.Status)
	if err != nil {
		return nil, err
	}

	result := make([]OrderDetails, 0, len(orders))
	for _, order := range orders {
		result = append(result, OrderDetails{
			OrderID:    order.ID,
			QuoteID:    order.QuoteID,
			CustomerID: order.CustomerID,
			Status:     order.Status,
			LineCount:  len(order.Lines),
		})
	}

	return result, nil
}
