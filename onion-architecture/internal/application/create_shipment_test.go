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
