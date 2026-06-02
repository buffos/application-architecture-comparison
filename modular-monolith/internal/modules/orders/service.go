package orders

import (
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/quotes"
)

type ApprovedQuoteSource interface {
	GetApprovedQuoteForOrder(quoteID string) (quotes.ApprovedQuote, error)
}

type ConvertQuoteToOrderCommand struct {
	QuoteID string
}

type ConvertQuoteToOrderResult struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type Service struct {
	orders    Repository
	quotes    ApprovedQuoteSource
	inventory inventory.Reserver
}

func NewService(orders Repository, quotes ApprovedQuoteSource, inventory inventory.Reserver) Service {
	return Service{
		orders:    orders,
		quotes:    quotes,
		inventory: inventory,
	}
}

func (s Service) ConvertQuoteToOrder(command ConvertQuoteToOrderCommand) (ConvertQuoteToOrderResult, error) {
	quote, err := s.quotes.GetApprovedQuoteForOrder(command.QuoteID)
	if err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	order := NewOrderFromApprovedQuote(quote)

	reservations := make([]inventory.ReservationItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		reservations = append(reservations, inventory.ReservationItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := s.inventory.Reserve(reservations); err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	return ConvertQuoteToOrderResult{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}
