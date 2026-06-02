package orders

import "modular-monolith/internal/modules/quotes"

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
	orders Repository
	quotes ApprovedQuoteSource
}

func NewService(orders Repository, quotes ApprovedQuoteSource) Service {
	return Service{
		orders: orders,
		quotes: quotes,
	}
}

func (s Service) ConvertQuoteToOrder(command ConvertQuoteToOrderCommand) (ConvertQuoteToOrderResult, error) {
	quote, err := s.quotes.GetApprovedQuoteForOrder(command.QuoteID)
	if err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	order := NewOrderFromApprovedQuote(quote)

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
