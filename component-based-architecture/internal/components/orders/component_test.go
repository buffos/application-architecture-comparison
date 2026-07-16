package orders

import (
	"errors"
	"testing"

	"component-based-architecture/internal/components/inventory"
	"component-based-architecture/internal/components/payments"
	"component-based-architecture/internal/components/quotes"
	"component-based-architecture/internal/components/shipments"
)

type stubApprovedQuoteSource struct {
	quote quotes.ApprovedQuote
	err   error
}

type stubReserver struct {
	items []inventory.ReservationItem
	err   error
}

func (s *stubReserver) Reserve(items []inventory.ReservationItem) error {
	s.items = append([]inventory.ReservationItem(nil), items...)
	return s.err
}

type stubPaymentProcessor struct {
	request payments.PaymentRequest
	err     error
}

func (s *stubPaymentProcessor) Capture(request payments.PaymentRequest) (payments.CaptureResult, error) {
	s.request = request
	return payments.CaptureResult{}, s.err
}

type stubShipmentCreator struct {
	request shipments.ShipmentRequest
	err     error
}

func (s *stubShipmentCreator) Create(request shipments.ShipmentRequest) (shipments.Shipment, error) {
	s.request = request
	if s.err != nil {
		return shipments.Shipment{}, s.err
	}
	return shipments.Shipment{ID: "shipment-001", OrderID: request.OrderID, CustomerID: request.CustomerID, Lines: request.Lines}, nil
}

func (s stubApprovedQuoteSource) GetApprovedQuoteForOrder(quoteID string) (quotes.ApprovedQuote, error) {
	if s.err != nil {
		return quotes.ApprovedQuote{}, s.err
	}
	return s.quote, nil
}

func TestConvertQuoteToOrderReservesStockAndCreatesOrderSnapshot(t *testing.T) {
	reserver := &stubReserver{}
	component := NewComponent(stubApprovedQuoteSource{quote: quotes.ApprovedQuote{
		QuoteID: "quote-001", CustomerID: "customer-001",
		Lines: []quotes.ApprovedQuoteLine{{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000}},
	}}, reserver, &stubPaymentProcessor{}, &stubShipmentCreator{})

	result, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("convert quote: %v", err)
	}
	if result.OrderID != "order-001" {
		t.Fatalf("expected order-001, got %s", result.OrderID)
	}
	if result.Status != OrderStatusPendingPayment {
		t.Fatalf("expected %s, got %s", OrderStatusPendingPayment, result.Status)
	}
	if result.LineCount != 1 {
		t.Fatalf("expected one line, got %d", result.LineCount)
	}
	if len(reserver.items) != 1 || reserver.items[0].ProductSKU != "sku-001" || reserver.items[0].Quantity != 2 {
		t.Fatalf("expected reservation for two sku-001 items, got %+v", reserver.items)
	}
}

func TestConvertQuoteToOrderPropagatesNonApprovedQuoteError(t *testing.T) {
	component := NewComponent(stubApprovedQuoteSource{err: quotes.ErrQuoteNotConvertible}, &stubReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{})

	_, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if !errors.Is(err, quotes.ErrQuoteNotConvertible) {
		t.Fatalf("expected %v, got %v", quotes.ErrQuoteNotConvertible, err)
	}
}

func TestConvertQuoteToOrderStopsWhenReservationFails(t *testing.T) {
	component := NewComponent(stubApprovedQuoteSource{quote: quotes.ApprovedQuote{
		QuoteID: "quote-001", CustomerID: "customer-001",
		Lines: []quotes.ApprovedQuoteLine{{ProductSKU: "sku-001", Quantity: 1}},
	}}, &stubReserver{err: inventory.ErrInsufficientStock}, &stubPaymentProcessor{}, &stubShipmentCreator{})

	_, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if !errors.Is(err, inventory.ErrInsufficientStock) {
		t.Fatalf("expected %v, got %v", inventory.ErrInsufficientStock, err)
	}
}

func TestCapturePaymentPaysPendingOrder(t *testing.T) {
	processor := &stubPaymentProcessor{}
	component := NewComponent(stubApprovedQuoteSource{quote: quotes.ApprovedQuote{
		QuoteID: "quote-001", CustomerID: "customer-001",
		Lines: []quotes.ApprovedQuoteLine{{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000}},
	}}, &stubReserver{}, processor, &stubShipmentCreator{})
	converted, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("convert quote: %v", err)
	}

	result, err := component.CapturePayment(CapturePaymentCommand{OrderID: converted.OrderID})
	if err != nil {
		t.Fatalf("capture payment: %v", err)
	}
	if result.Status != OrderStatusPaid {
		t.Fatalf("expected %s, got %s", OrderStatusPaid, result.Status)
	}
	if processor.request.Amount != 30000 {
		t.Fatalf("expected capture amount 30000, got %d", processor.request.Amount)
	}
}

func TestCapturePaymentRejectsAlreadyPaidOrder(t *testing.T) {
	component := NewComponent(stubApprovedQuoteSource{quote: quotes.ApprovedQuote{
		QuoteID: "quote-001", CustomerID: "customer-001", Lines: []quotes.ApprovedQuoteLine{{ProductSKU: "sku-001", Quantity: 1, UnitPrice: 15000}},
	}}, &stubReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{})
	converted, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("convert quote: %v", err)
	}
	if _, err := component.CapturePayment(CapturePaymentCommand{OrderID: converted.OrderID}); err != nil {
		t.Fatalf("capture payment: %v", err)
	}

	_, err = component.CapturePayment(CapturePaymentCommand{OrderID: converted.OrderID})
	if !errors.Is(err, ErrOrderNotPayable) {
		t.Fatalf("expected %v, got %v", ErrOrderNotPayable, err)
	}
}

func TestCreateShipmentShipsPaidOrder(t *testing.T) {
	creator := &stubShipmentCreator{}
	component := NewComponent(stubApprovedQuoteSource{quote: quotes.ApprovedQuote{
		QuoteID: "quote-001", CustomerID: "customer-001", Lines: []quotes.ApprovedQuoteLine{{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 1, UnitPrice: 15000}},
	}}, &stubReserver{}, &stubPaymentProcessor{}, creator)
	converted, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("convert quote: %v", err)
	}
	if _, err := component.CapturePayment(CapturePaymentCommand{OrderID: converted.OrderID}); err != nil {
		t.Fatalf("capture payment: %v", err)
	}

	result, err := component.CreateShipment(CreateShipmentCommand{OrderID: converted.OrderID})
	if err != nil {
		t.Fatalf("create shipment: %v", err)
	}
	if result.Status != OrderStatusShipped {
		t.Fatalf("expected %s, got %s", OrderStatusShipped, result.Status)
	}
	if creator.request.OrderID != converted.OrderID || len(creator.request.Lines) != 1 {
		t.Fatalf("expected shipment request for order, got %+v", creator.request)
	}
}

func TestCreateShipmentRejectsUnpaidOrder(t *testing.T) {
	component := NewComponent(stubApprovedQuoteSource{quote: quotes.ApprovedQuote{
		QuoteID: "quote-001", CustomerID: "customer-001", Lines: []quotes.ApprovedQuoteLine{{ProductSKU: "sku-001", Quantity: 1}},
	}}, &stubReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{})
	converted, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("convert quote: %v", err)
	}

	_, err = component.CreateShipment(CreateShipmentCommand{OrderID: converted.OrderID})
	if !errors.Is(err, ErrOrderNotShippable) {
		t.Fatalf("expected %v, got %v", ErrOrderNotShippable, err)
	}
}
