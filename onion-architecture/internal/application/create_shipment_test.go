package application

import (
	"testing"
	"time"

	"onion-architecture/internal/domain"
	timeinfra "onion-architecture/internal/infrastructure/services/time"
)

type stubShipmentStore struct {
	saved domain.Shipment
}

func (s *stubShipmentStore) Save(shipment domain.Shipment) error {
	s.saved = shipment
	return nil
}

func TestCreateShipmentServiceShipsPaidOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPaid,
			Lines: []domain.OrderLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        1,
					UnitPrice:       45000,
				},
			},
		},
	}
	shipments := &stubShipmentStore{}
	clock := timeinfra.NewFixedClock(time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))

	service := NewCreateShipmentService(orders, shipments, clock)

	result, err := service.Execute(CreateShipmentCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.OrderStatus != domain.OrderStatusShipped {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusShipped, result.OrderStatus)
	}

	if shipments.saved.OrderID != "order-001" {
		t.Fatalf("expected shipment order id order-001, got %s", shipments.saved.OrderID)
	}

	if orders.saved.Status != domain.OrderStatusShipped {
		t.Fatalf("expected saved order status %s, got %s", domain.OrderStatusShipped, orders.saved.Status)
	}

	if !orders.saved.ShippedAt.Equal(clock.Now()) {
		t.Fatalf("expected shipped at %v, got %v", clock.Now(), orders.saved.ShippedAt)
	}
}

func TestCreateShipmentServiceRejectsUnpaidOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPendingPayment,
		},
	}
	shipments := &stubShipmentStore{}
	clock := timeinfra.NewFixedClock(time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))

	service := NewCreateShipmentService(orders, shipments, clock)

	_, err := service.Execute(CreateShipmentCommand{OrderID: "order-001"})
	if err != domain.ErrOrderNotShippable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotShippable, err)
	}
}

func TestCreateShipmentServiceRejectsOrderInPaymentReview(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPaymentReview,
		},
	}
	shipments := &stubShipmentStore{}
	clock := timeinfra.NewFixedClock(time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))

	service := NewCreateShipmentService(orders, shipments, clock)

	_, err := service.Execute(CreateShipmentCommand{OrderID: "order-001"})
	if err != domain.ErrOrderNotShippable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotShippable, err)
	}
}

func TestCreateShipmentServiceSupportsPartialShipment(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPaid,
			Lines: []domain.OrderLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        2,
					UnitPrice:       45000,
				},
			},
		},
	}
	shipments := &stubShipmentStore{}
	clock := timeinfra.NewFixedClock(time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC))

	service := NewCreateShipmentService(orders, shipments, clock)

	result, err := service.Execute(CreateShipmentCommand{
		OrderID: "order-001",
		Lines: []CreateShipmentLine{
			{ProductSKU: "sku-002", Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.OrderStatus != domain.OrderStatusPartiallyShipped {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusPartiallyShipped, result.OrderStatus)
	}

	if len(shipments.saved.Lines) != 1 || shipments.saved.Lines[0].Quantity != 1 {
		t.Fatalf("expected one shipment line with quantity 1, got %+v", shipments.saved.Lines)
	}

	if orders.saved.Lines[0].ShippedQuantity != 1 {
		t.Fatalf("expected shipped quantity 1, got %d", orders.saved.Lines[0].ShippedQuantity)
	}

	if !orders.saved.ShippedAt.Equal(clock.Now()) {
		t.Fatalf("expected shipped at %v, got %v", clock.Now(), orders.saved.ShippedAt)
	}
}

func TestCreateShipmentServiceShipsRemainingQuantityAfterPartialShipment(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPartiallyShipped,
			ShippedAt:  time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        2,
					ShippedQuantity: 1,
					UnitPrice:       45000,
				},
			},
		},
	}
	shipments := &stubShipmentStore{}
	clock := timeinfra.NewFixedClock(time.Date(2026, 6, 2, 10, 0, 0, 0, time.UTC))

	service := NewCreateShipmentService(orders, shipments, clock)

	result, err := service.Execute(CreateShipmentCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.OrderStatus != domain.OrderStatusShipped {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusShipped, result.OrderStatus)
	}

	if len(shipments.saved.Lines) != 1 || shipments.saved.Lines[0].Quantity != 1 {
		t.Fatalf("expected remaining shipment quantity 1, got %+v", shipments.saved.Lines)
	}

	if orders.saved.Lines[0].ShippedQuantity != 2 {
		t.Fatalf("expected shipped quantity 2, got %d", orders.saved.Lines[0].ShippedQuantity)
	}

	if !orders.saved.ShippedAt.Equal(time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected first shipped at timestamp to be preserved, got %v", orders.saved.ShippedAt)
	}
}
