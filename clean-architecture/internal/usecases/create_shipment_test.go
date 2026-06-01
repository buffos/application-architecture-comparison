package usecases

import (
	"testing"
	"time"

	"clean-architecture/internal/entities"
)

type stubShipmentWriter struct {
	saved entities.Shipment
}

func (g *stubShipmentWriter) Save(shipment entities.Shipment) error {
	g.saved = shipment
	return nil
}

type stubCreateShipmentOutput struct {
	output CreateShipmentOutput
}

func (o *stubCreateShipmentOutput) Present(output CreateShipmentOutput) error {
	o.output = output
	return nil
}

func TestCreateShipmentInteractorCreatesShipmentForPaidOrder(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPaid,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", ProductName: "Office Chair", Quantity: 2},
			},
		},
	}
	shipments := &stubShipmentWriter{}
	output := &stubCreateShipmentOutput{}
	clock := stubClock{now: time.Date(2026, 6, 5, 9, 0, 0, 0, time.UTC)}

	interactor := NewCreateShipmentInteractor(orders, shipments, clock, output)

	err := interactor.Execute(CreateShipmentInput{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if shipments.saved.OrderID != "order-001" {
		t.Fatalf("expected shipment order id order-001, got %s", shipments.saved.OrderID)
	}

	if orders.saved.Status != entities.OrderStatusShipped {
		t.Fatalf("expected order status %s, got %s", entities.OrderStatusShipped, orders.saved.Status)
	}
	if orders.saved.Lines[0].ShippedQuantity != 2 {
		t.Fatalf("expected shipped quantity 2, got %d", orders.saved.Lines[0].ShippedQuantity)
	}
	if orders.saved.ShippedAt == nil {
		t.Fatal("expected shipped timestamp to be set")
	}

	if output.output.ShipmentID == "" {
		t.Fatal("expected presenter output to include shipment id")
	}
}

func TestCreateShipmentInteractorRejectsUnpaidOrder(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-002",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPendingPayment,
		},
	}
	shipments := &stubShipmentWriter{}
	output := &stubCreateShipmentOutput{}
	clock := stubClock{now: time.Date(2026, 6, 5, 9, 0, 0, 0, time.UTC)}

	interactor := NewCreateShipmentInteractor(orders, shipments, clock, output)

	err := interactor.Execute(CreateShipmentInput{OrderID: "order-002"})
	if err != entities.ErrQuoteCannotTransition {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotTransition, err)
	}
}

func TestCreateShipmentInteractorCreatesPartialShipment(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-004",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPaid,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", ProductName: "Office Chair", Quantity: 5},
			},
		},
	}
	shipments := &stubShipmentWriter{}
	output := &stubCreateShipmentOutput{}
	clock := stubClock{now: time.Date(2026, 6, 5, 9, 0, 0, 0, time.UTC)}

	interactor := NewCreateShipmentInteractor(orders, shipments, clock, output)

	err := interactor.Execute(CreateShipmentInput{
		OrderID: "order-004",
		Lines: []CreateShipmentLineInput{
			{SKU: "CHAIR-001", Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if shipments.saved.Lines[0].Quantity != 2 {
		t.Fatalf("expected shipped quantity 2, got %d", shipments.saved.Lines[0].Quantity)
	}

	if orders.saved.Status != entities.OrderStatusPartiallyShipped {
		t.Fatalf("expected order status %s, got %s", entities.OrderStatusPartiallyShipped, orders.saved.Status)
	}

	if orders.saved.Lines[0].ShippedQuantity != 2 {
		t.Fatalf("expected shipped quantity 2, got %d", orders.saved.Lines[0].ShippedQuantity)
	}
}

func TestCreateShipmentInteractorAllowsShippingRemainingLinesAfterPartial(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-005",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPartiallyShipped,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", ProductName: "Office Chair", Quantity: 5, ShippedQuantity: 2},
			},
		},
	}
	shipments := &stubShipmentWriter{}
	output := &stubCreateShipmentOutput{}
	clock := stubClock{now: time.Date(2026, 6, 6, 9, 0, 0, 0, time.UTC)}

	interactor := NewCreateShipmentInteractor(orders, shipments, clock, output)

	err := interactor.Execute(CreateShipmentInput{OrderID: "order-005"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if shipments.saved.Lines[0].Quantity != 3 {
		t.Fatalf("expected shipped quantity 3, got %d", shipments.saved.Lines[0].Quantity)
	}

	if orders.saved.Status != entities.OrderStatusShipped {
		t.Fatalf("expected order status %s, got %s", entities.OrderStatusShipped, orders.saved.Status)
	}

	if orders.saved.Lines[0].ShippedQuantity != 5 {
		t.Fatalf("expected shipped quantity 5, got %d", orders.saved.Lines[0].ShippedQuantity)
	}
}

func TestCreateShipmentInteractorRejectsOrderInPaymentReview(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-003",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPaymentReview,
		},
	}
	shipments := &stubShipmentWriter{}
	output := &stubCreateShipmentOutput{}
	clock := stubClock{now: time.Date(2026, 6, 5, 9, 0, 0, 0, time.UTC)}

	interactor := NewCreateShipmentInteractor(orders, shipments, clock, output)

	err := interactor.Execute(CreateShipmentInput{OrderID: "order-003"})
	if err != entities.ErrQuoteCannotTransition {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotTransition, err)
	}
}
