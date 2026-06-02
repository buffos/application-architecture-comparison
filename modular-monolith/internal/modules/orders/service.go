package orders

import (
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/payments"
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

type CapturePaymentCommand struct {
	OrderID string
}

type CapturePaymentResult struct {
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type Service struct {
	orders    Repository
	quotes    ApprovedQuoteSource
	inventory inventory.Reserver
	payments  payments.Processor
}

func NewService(orders Repository, quotes ApprovedQuoteSource, inventory inventory.Reserver, payments payments.Processor) Service {
	return Service{
		orders:    orders,
		quotes:    quotes,
		inventory: inventory,
		payments:  payments,
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

func (s Service) CapturePayment(command CapturePaymentCommand) (CapturePaymentResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return CapturePaymentResult{}, err
	}

	totalAmount := 0
	for _, line := range order.Lines {
		totalAmount += line.Quantity * line.UnitPrice
	}

	if err := s.payments.Capture(payments.PaymentRequest{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Amount:     totalAmount,
	}); err != nil {
		return CapturePaymentResult{}, err
	}

	if err := order.MarkPaid(); err != nil {
		return CapturePaymentResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return CapturePaymentResult{}, err
	}

	return CapturePaymentResult{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}
