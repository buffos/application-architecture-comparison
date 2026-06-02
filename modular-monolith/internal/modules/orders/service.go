package orders

import (
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/payments"
	"modular-monolith/internal/modules/quotes"
	"modular-monolith/internal/modules/shipments"
	"time"
)

type ApprovedQuoteSource interface {
	GetApprovedQuoteForOrder(quoteID string) (quotes.ApprovedQuote, error)
}

type Clock interface {
	Now() time.Time
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

type ApprovePaymentReviewCommand struct {
	OrderID string
}

type ApprovePaymentReviewResult struct {
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type CreateShipmentCommand struct {
	OrderID string
}

type CreateShipmentResult struct {
	ShipmentID string
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type ReturnableOrder struct {
	OrderID    string
	CustomerID string
	ShippedAt  time.Time
	Lines      []ReturnableOrderLine
}

type ReturnableOrderLine struct {
	ProductSKU       string
	ProductName      string
	ProductCategory  string
	Quantity         int
	UnitPrice        int
	ReturnWindowDays int
}

type CancelOrderCommand struct {
	OrderID string
}

type CancelOrderResult struct {
	OrderID    string
	CustomerID string
	Status     string
	LineCount  int
}

type Service struct {
	orders    Repository
	quotes    ApprovedQuoteSource
	inventory inventory.StockKeeper
	payments  payments.Processor
	shipments shipments.Creator
	clock     Clock
}

func NewService(orders Repository, quotes ApprovedQuoteSource, inventory inventory.StockKeeper, payments payments.Processor, shipments shipments.Creator, clock Clock) Service {
	return Service{
		orders:    orders,
		quotes:    quotes,
		inventory: inventory,
		payments:  payments,
		shipments: shipments,
		clock:     clock,
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

	result, err := s.payments.Capture(payments.PaymentRequest{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Amount:     totalAmount,
	})
	if err != nil {
		return CapturePaymentResult{}, err
	}

	switch result.Outcome {
	case payments.CaptureOutcomeReview:
		if err := order.MarkPaymentReview(); err != nil {
			return CapturePaymentResult{}, err
		}
	default:
		if err := order.MarkPaid(); err != nil {
			return CapturePaymentResult{}, err
		}
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

func (s Service) ApprovePaymentReview(command ApprovePaymentReviewCommand) (ApprovePaymentReviewResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return ApprovePaymentReviewResult{}, err
	}

	if err := order.ApprovePaymentReview(); err != nil {
		return ApprovePaymentReviewResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return ApprovePaymentReviewResult{}, err
	}

	return ApprovePaymentReviewResult{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}

func (s Service) CreateShipment(command CreateShipmentCommand) (CreateShipmentResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return CreateShipmentResult{}, err
	}

	shipmentLines := make([]shipments.ShipmentLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		shipmentLines = append(shipmentLines, shipments.ShipmentLine{
			ProductSKU:  line.ProductSKU,
			ProductName: line.ProductName,
			Quantity:    line.Quantity,
		})
	}

	shipment, err := s.shipments.Create(shipments.ShipmentRequest{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Lines:      shipmentLines,
	})
	if err != nil {
		return CreateShipmentResult{}, err
	}

	if err := order.MarkShipped(s.clock.Now()); err != nil {
		return CreateShipmentResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return CreateShipmentResult{}, err
	}

	return CreateShipmentResult{
		ShipmentID: shipment.ID,
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}

func (s Service) CancelOrder(command CancelOrderCommand) (CancelOrderResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return CancelOrderResult{}, err
	}

	if err := order.Cancel(); err != nil {
		return CancelOrderResult{}, err
	}

	releases := make([]inventory.ReleaseItem, 0, len(order.Lines))
	for _, line := range order.Lines {
		releases = append(releases, inventory.ReleaseItem{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	if err := s.inventory.Release(releases); err != nil {
		return CancelOrderResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return CancelOrderResult{}, err
	}

	return CancelOrderResult{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}

func (s Service) GetReturnableOrder(orderID string) (ReturnableOrder, error) {
	order, err := s.orders.FindByID(orderID)
	if err != nil {
		return ReturnableOrder{}, err
	}

	if err := order.EnsureReturnable(); err != nil {
		return ReturnableOrder{}, err
	}

	lines := make([]ReturnableOrderLine, 0, len(order.Lines))
	for _, line := range order.Lines {
		lines = append(lines, ReturnableOrderLine{
			ProductSKU:       line.ProductSKU,
			ProductName:      line.ProductName,
			ProductCategory:  line.ProductCategory,
			Quantity:         line.Quantity,
			UnitPrice:        line.UnitPrice,
			ReturnWindowDays: line.ReturnWindowDays,
		})
	}

	return ReturnableOrder{
		OrderID:    order.ID,
		CustomerID: order.CustomerID,
		ShippedAt:  order.ShippedAt,
		Lines:      lines,
	}, nil
}
