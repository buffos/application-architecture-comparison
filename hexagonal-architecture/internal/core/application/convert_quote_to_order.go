package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ConvertQuoteToOrderUseCase struct {
	quotes    ports.QuoteRepository
	orders    ports.OrderRepository
	inventory ports.InventoryReservation
}

func NewConvertQuoteToOrderUseCase(quotes ports.QuoteRepository, orders ports.OrderRepository, inventory ports.InventoryReservation) ConvertQuoteToOrderUseCase {
	return ConvertQuoteToOrderUseCase{
		quotes:    quotes,
		orders:    orders,
		inventory: inventory,
	}
}

func (uc ConvertQuoteToOrderUseCase) Execute(id string) (domain.Order, error) {
	quote, err := uc.quotes.FindByID(id)
	if err != nil {
		return domain.Order{}, err
	}

	reservations := make([]domain.ReservationLine, 0, len(quote.Lines))
	for _, line := range quote.Lines {
		reservations = append(reservations, domain.ReservationLine{
			SKU:      line.SKU,
			Quantity: line.Quantity,
		})
	}

	if err := uc.inventory.Reserve(reservations); err != nil {
		return domain.Order{}, err
	}

	order, err := domain.NewOrderFromQuote(quote)
	if err != nil {
		return domain.Order{}, err
	}

	if err := uc.orders.Save(order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
