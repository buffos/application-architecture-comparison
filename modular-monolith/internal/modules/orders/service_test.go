package orders

import (
	"testing"
	"time"

	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/payments"
	"modular-monolith/internal/modules/quotes"
	"modular-monolith/internal/modules/shipments"
)

type stubOrderRepository struct {
	saved Order
}

func (r *stubOrderRepository) Save(order Order) error {
	r.saved = order
	return nil
}

func (r *stubOrderRepository) FindByID(id string) (Order, error) {
	return r.saved, nil
}

func (r *stubOrderRepository) ListByStatus(status string) ([]Order, error) {
	if r.saved.ID == "" {
		return nil, nil
	}
	if status == "" || r.saved.Status == status {
		return []Order{r.saved}, nil
	}
	return nil, nil
}

type stubApprovedQuoteSource struct {
	quote quotes.ApprovedQuote
	err   error
}

type stubInventoryReserver struct {
	reserved []inventory.ReservationItem
	released []inventory.ReleaseItem
	err      error
}

type stubPaymentProcessor struct {
	request PaymentRequestAlias
	result  payments.CaptureResult
	err     error
}

type PaymentRequestAlias = payments.PaymentRequest

type stubShipmentCreator struct {
	request  shipments.ShipmentRequest
	shipment shipments.Shipment
	err      error
}

type stubClock struct {
	now time.Time
}

func (s *stubInventoryReserver) Reserve(items []inventory.ReservationItem) error {
	if s.err != nil {
		return s.err
	}

	s.reserved = append([]inventory.ReservationItem(nil), items...)
	return nil
}

func (s *stubInventoryReserver) Release(items []inventory.ReleaseItem) error {
	if s.err != nil {
		return s.err
	}

	s.released = append([]inventory.ReleaseItem(nil), items...)
	return nil
}

func (s *stubPaymentProcessor) Capture(request payments.PaymentRequest) (payments.CaptureResult, error) {
	if s.err != nil {
		return payments.CaptureResult{}, s.err
	}

	s.request = request
	if s.result.Outcome == "" {
		s.result = payments.CaptureResult{Outcome: payments.CaptureOutcomeApproved}
	}
	return s.result, nil
}

func (s *stubShipmentCreator) Create(request shipments.ShipmentRequest) (shipments.Shipment, error) {
	if s.err != nil {
		return shipments.Shipment{}, s.err
	}

	s.request = request
	if s.shipment.ID == "" {
		s.shipment = shipments.Shipment{
			ID:         "shipment-001",
			OrderID:    request.OrderID,
			CustomerID: request.CustomerID,
			Lines:      request.Lines,
		}
	}

	return s.shipment, nil
}

func (c stubClock) Now() time.Time {
	return c.now
}

func (s stubApprovedQuoteSource) GetApprovedQuoteForOrder(quoteID string) (quotes.ApprovedQuote, error) {
	if s.err != nil {
		return quotes.ApprovedQuote{}, s.err
	}

	return s.quote, nil
}

func TestConvertQuoteToOrderCreatesPendingPaymentOrder(t *testing.T) {
	orders := &stubOrderRepository{}
	reserver := &stubInventoryReserver{}
	paymentProcessor := &stubPaymentProcessor{}
	shipmentCreator := &stubShipmentCreator{}
	service := NewService(orders, stubApprovedQuoteSource{
		quote: quotes.ApprovedQuote{
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Lines: []quotes.ApprovedQuoteLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
			},
		},
	}, reserver, paymentProcessor, shipmentCreator, stubClock{})

	result, err := service.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusPendingPayment {
		t.Fatalf("expected status %s, got %s", OrderStatusPendingPayment, result.Status)
	}

	if orders.saved.QuoteID != "quote-001" {
		t.Fatalf("expected quote-001, got %s", orders.saved.QuoteID)
	}

	if len(reserver.reserved) != 1 || reserver.reserved[0].Quantity != 2 {
		t.Fatalf("expected reservation for quantity 2, got %+v", reserver.reserved)
	}
}

func TestConvertQuoteToOrderRejectsNonApprovedQuote(t *testing.T) {
	orders := &stubOrderRepository{}
	reserver := &stubInventoryReserver{}
	paymentProcessor := &stubPaymentProcessor{}
	shipmentCreator := &stubShipmentCreator{}
	service := NewService(orders, stubApprovedQuoteSource{
		err: quotes.ErrQuoteNotConvertible,
	}, reserver, paymentProcessor, shipmentCreator, stubClock{})

	_, err := service.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != quotes.ErrQuoteNotConvertible {
		t.Fatalf("expected %v, got %v", quotes.ErrQuoteNotConvertible, err)
	}
}

func TestConvertQuoteToOrderStopsWhenReservationFails(t *testing.T) {
	orders := &stubOrderRepository{}
	reserver := &stubInventoryReserver{err: inventory.ErrInsufficientStock}
	paymentProcessor := &stubPaymentProcessor{}
	shipmentCreator := &stubShipmentCreator{}
	service := NewService(orders, stubApprovedQuoteSource{
		quote: quotes.ApprovedQuote{
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Lines: []quotes.ApprovedQuoteLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
			},
		},
	}, reserver, paymentProcessor, shipmentCreator, stubClock{})

	_, err := service.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != inventory.ErrInsufficientStock {
		t.Fatalf("expected %v, got %v", inventory.ErrInsufficientStock, err)
	}

	if orders.saved.ID != "" {
		t.Fatalf("expected order not to be saved when reservation fails")
	}
}

func TestCapturePaymentMarksOrderPaid(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:         "order-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPendingPayment,
			Lines: []OrderLine{
				{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000},
			},
		},
	}
	paymentProcessor := &stubPaymentProcessor{}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, paymentProcessor, &stubShipmentCreator{}, stubClock{})

	result, err := service.CapturePayment(CapturePaymentCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusPaid {
		t.Fatalf("expected %s, got %s", OrderStatusPaid, result.Status)
	}

	if paymentProcessor.request.Amount != 30000 {
		t.Fatalf("expected amount 30000, got %d", paymentProcessor.request.Amount)
	}
}

func TestCapturePaymentMovesOrderToPaymentReviewWhenGatewayRequestsReview(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:         "order-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPendingPayment,
			Lines: []OrderLine{
				{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000},
			},
		},
	}
	paymentProcessor := &stubPaymentProcessor{
		result: payments.CaptureResult{Outcome: payments.CaptureOutcomeReview},
	}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, paymentProcessor, &stubShipmentCreator{}, stubClock{})

	result, err := service.CapturePayment(CapturePaymentCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusPaymentReview {
		t.Fatalf("expected %s, got %s", OrderStatusPaymentReview, result.Status)
	}
}

func TestCapturePaymentRejectsOrderThatIsNotPayable(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:     "order-001",
			Status: OrderStatusPaid,
		},
	}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{}, stubClock{})

	_, err := service.CapturePayment(CapturePaymentCommand{OrderID: "order-001"})
	if err != ErrOrderNotPayable {
		t.Fatalf("expected %v, got %v", ErrOrderNotPayable, err)
	}
}

func TestCreateShipmentMarksOrderShipped(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:         "order-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPaid,
			Lines: []OrderLine{
				{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 2},
			},
		},
	}
	shipmentCreator := &stubShipmentCreator{}
	clock := stubClock{now: time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, shipmentCreator, clock)

	result, err := service.CreateShipment(CreateShipmentCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusShipped {
		t.Fatalf("expected %s, got %s", OrderStatusShipped, result.Status)
	}

	if shipmentCreator.request.OrderID != "order-001" {
		t.Fatalf("expected shipment for order-001, got %s", shipmentCreator.request.OrderID)
	}

	if !orders.saved.ShippedAt.Equal(clock.now) {
		t.Fatalf("expected shipped time to be recorded")
	}
}

func TestCreateShipmentSupportsPartialShipment(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:         "order-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPaid,
			Lines: []OrderLine{
				{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 3},
			},
		},
	}
	shipmentCreator := &stubShipmentCreator{}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, shipmentCreator, stubClock{})

	result, err := service.CreateShipment(CreateShipmentCommand{
		OrderID: "order-001",
		Lines: []CreateShipmentLine{
			{ProductSKU: "sku-001", Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusPartiallyShipped {
		t.Fatalf("expected %s, got %s", OrderStatusPartiallyShipped, result.Status)
	}

	if len(shipmentCreator.request.Lines) != 1 || shipmentCreator.request.Lines[0].Quantity != 1 {
		t.Fatalf("expected one shipped unit, got %+v", shipmentCreator.request.Lines)
	}

	if orders.saved.Lines[0].ShippedQuantity != 1 {
		t.Fatalf("expected shipped quantity 1, got %d", orders.saved.Lines[0].ShippedQuantity)
	}
}

func TestCreateShipmentCanShipRemainingQuantityLater(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:         "order-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPartiallyShipped,
			Lines: []OrderLine{
				{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 3, ShippedQuantity: 1},
			},
		},
	}
	shipmentCreator := &stubShipmentCreator{}
	clock := stubClock{now: time.Date(2026, 6, 2, 12, 0, 0, 0, time.UTC)}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, shipmentCreator, clock)

	result, err := service.CreateShipment(CreateShipmentCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusShipped {
		t.Fatalf("expected %s, got %s", OrderStatusShipped, result.Status)
	}

	if len(shipmentCreator.request.Lines) != 1 || shipmentCreator.request.Lines[0].Quantity != 2 {
		t.Fatalf("expected remaining quantity 2, got %+v", shipmentCreator.request.Lines)
	}

	if orders.saved.Lines[0].ShippedQuantity != 3 {
		t.Fatalf("expected shipped quantity 3, got %d", orders.saved.Lines[0].ShippedQuantity)
	}

	if !orders.saved.ShippedAt.Equal(clock.now) {
		t.Fatalf("expected shipped time to be recorded on final shipment")
	}
}

func TestCreateShipmentRejectsOrderThatIsNotPaid(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:     "order-001",
			Status: OrderStatusPendingPayment,
		},
	}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{}, stubClock{})

	_, err := service.CreateShipment(CreateShipmentCommand{OrderID: "order-001"})
	if err != ErrOrderNotShippable {
		t.Fatalf("expected %v, got %v", ErrOrderNotShippable, err)
	}
}

func TestCreateShipmentRejectsOrderThatIsInPaymentReview(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:     "order-001",
			Status: OrderStatusPaymentReview,
		},
	}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{}, stubClock{})

	_, err := service.CreateShipment(CreateShipmentCommand{OrderID: "order-001"})
	if err != ErrOrderNotShippable {
		t.Fatalf("expected %v, got %v", ErrOrderNotShippable, err)
	}
}

func TestApprovePaymentReviewMovesOrderToPaid(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:         "order-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPaymentReview,
			Lines: []OrderLine{
				{ProductSKU: "sku-001", Quantity: 2, UnitPrice: 15000},
			},
		},
	}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{}, stubClock{})

	result, err := service.ApprovePaymentReview(ApprovePaymentReviewCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusPaid {
		t.Fatalf("expected %s, got %s", OrderStatusPaid, result.Status)
	}
}

func TestCancelOrderReleasesInventoryAndMarksOrderCancelled(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:         "order-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPendingPayment,
			Lines: []OrderLine{
				{ProductSKU: "sku-001", Quantity: 2},
			},
		},
	}
	inventoryModule := &stubInventoryReserver{}
	service := NewService(orders, stubApprovedQuoteSource{}, inventoryModule, &stubPaymentProcessor{}, &stubShipmentCreator{}, stubClock{})

	result, err := service.CancelOrder(CancelOrderCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusCancelled {
		t.Fatalf("expected %s, got %s", OrderStatusCancelled, result.Status)
	}

	if len(inventoryModule.released) != 1 || inventoryModule.released[0].Quantity != 2 {
		t.Fatalf("expected release of quantity 2, got %+v", inventoryModule.released)
	}
}

func TestCancelOrderRejectsShippedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:     "order-001",
			Status: OrderStatusShipped,
		},
	}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{}, stubClock{})

	_, err := service.CancelOrder(CancelOrderCommand{OrderID: "order-001"})
	if err != ErrOrderNotCancellable {
		t.Fatalf("expected %v, got %v", ErrOrderNotCancellable, err)
	}
}

func TestCancelOrderRejectsPartiallyShippedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		saved: Order{
			ID:     "order-001",
			Status: OrderStatusPartiallyShipped,
		},
	}
	service := NewService(orders, stubApprovedQuoteSource{}, &stubInventoryReserver{}, &stubPaymentProcessor{}, &stubShipmentCreator{}, stubClock{})

	_, err := service.CancelOrder(CancelOrderCommand{OrderID: "order-001"})
	if err != ErrOrderNotCancellable {
		t.Fatalf("expected %v, got %v", ErrOrderNotCancellable, err)
	}
}
