package application

import "onion-architecture/internal/domain"

type OrderFinder interface {
	FindByID(id string) (domain.Order, error)
	ListByStatus(status string) ([]domain.Order, error)
}

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

type GetOrderService struct {
	orders OrderFinder
}

func NewGetOrderService(orders OrderFinder) GetOrderService {
	return GetOrderService{
		orders: orders,
	}
}

func (s GetOrderService) Execute(query GetOrderQuery) (OrderDetails, error) {
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
