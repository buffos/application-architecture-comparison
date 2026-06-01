package application

import "onion-architecture/internal/domain"

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

type OrderStore interface {
	Save(order domain.Order) error
}

type ConvertQuoteToOrderService struct {
	quotes QuoteFinder
	orders OrderStore
}

func NewConvertQuoteToOrderService(quotes QuoteFinder, orders OrderStore) ConvertQuoteToOrderService {
	return ConvertQuoteToOrderService{
		quotes: quotes,
		orders: orders,
	}
}

func (s ConvertQuoteToOrderService) Execute(command ConvertQuoteToOrderCommand) (ConvertQuoteToOrderResult, error) {
	quote, err := s.quotes.FindByID(command.QuoteID)
	if err != nil {
		return ConvertQuoteToOrderResult{}, err
	}

	order, err := domain.NewOrderFromQuote(quote)
	if err != nil {
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
