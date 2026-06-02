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
	Lines      []OrderLineDetails
}

type ListOrdersQuery struct {
	Status string
}

type OrderLineDetails struct {
	ProductSKU      string
	ProductCategory string
	Quantity        int
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
		Lines:      toOrderLineDetails(order.Lines),
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
			Lines:      toOrderLineDetails(order.Lines),
		})
	}

	return list, nil
}

func toOrderLineDetails(lines []OrderLine) []OrderLineDetails {
	details := make([]OrderLineDetails, 0, len(lines))
	for _, line := range lines {
		details = append(details, OrderLineDetails{
			ProductSKU:      line.ProductSKU,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
		})
	}

	return details
}
