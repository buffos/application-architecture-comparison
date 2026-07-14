package orders

import "microkernel-architecture/internal/kernel"

type Service struct {
	orders Repository
	quotes kernel.ApprovedQuoteProvider
	stock  kernel.InventoryReservation
}

func NewService(orders Repository, quotes kernel.ApprovedQuoteProvider, stock kernel.InventoryReservation) Service {
	return Service{
		orders: orders,
		quotes: quotes,
		stock:  stock,
	}
}

func (s Service) ConvertQuoteToOrder(command kernel.ConvertQuoteToOrderCommand) (kernel.ConvertQuoteToOrderResult, error) {
	quote, err := s.quotes.GetApprovedQuoteForOrder(command.QuoteID)
	if err != nil {
		return kernel.ConvertQuoteToOrderResult{}, err
	}

	lines := make([]OrderLine, 0, len(quote.Lines))
	reservationItems := make([]kernel.InventoryReservationItem, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		lines = append(lines, OrderLine{
			ProductSKU:      line.ProductSKU,
			ProductName:     line.ProductName,
			ProductCategory: line.ProductCategory,
			Quantity:        line.Quantity,
			UnitPrice:       line.UnitPrice,
		})
		reservationItems = append(reservationItems, kernel.InventoryReservationItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	order := NewOrderFromApprovedQuote(quote.QuoteID, quote.CustomerID, lines)
	if err := s.stock.Reserve(reservationItems); err != nil {
		return kernel.ConvertQuoteToOrderResult{}, err
	}

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
