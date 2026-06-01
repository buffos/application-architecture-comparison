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

type InventoryReservation interface {
	Reserve(items []domain.InventoryReservationItem) error
}

type ConvertQuoteToOrderService struct {
	quotes QuoteFinder
	orders OrderStore
	inventory InventoryReservation
}

func NewConvertQuoteToOrderService(quotes QuoteFinder, orders OrderStore, inventory InventoryReservation) ConvertQuoteToOrderService {
	return ConvertQuoteToOrderService{
		quotes:    quotes,
		orders:    orders,
		inventory: inventory,
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

	items := make([]domain.InventoryReservationItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		items = append(items, domain.InventoryReservationItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := s.inventory.Reserve(items); err != nil {
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
