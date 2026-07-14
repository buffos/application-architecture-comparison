package orders

import "microkernel-architecture/internal/kernel"

type Service struct {
	orders Repository
	quotes kernel.ApprovedQuoteProvider
}

func NewService(orders Repository, quotes kernel.ApprovedQuoteProvider) Service {
	return Service{
		orders: orders,
		quotes: quotes,
	}
}

func (s Service) ConvertQuoteToOrder(command kernel.ConvertQuoteToOrderCommand) (kernel.ConvertQuoteToOrderResult, error) {
	quote, err := s.quotes.GetApprovedQuoteForOrder(command.QuoteID)
	if err != nil {
		return kernel.ConvertQuoteToOrderResult{}, err
	}

	lines := make([]OrderLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU:      line.ProductSKU,
			ProductName:     line.ProductName,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
			UnitPrice:       line.UnitPrice,
		})
	}

	order := NewOrderFromApprovedQuote(quote.QuoteID, quote.CustomerID, lines)
	if err := s.orders.Save(order); err != nil {
		return kernel.ConvertQuoteToOrderResult{}, err
	}

	return kernel.ConvertQuoteToOrderResult{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}
