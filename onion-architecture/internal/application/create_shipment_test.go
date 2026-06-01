package application

import (
	"testing"

	"onion-architecture/internal/domain"
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

	service := NewCreateShipmentService(orders, shipments)

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

	service := NewCreateShipmentService(orders, shipments)

	_, err := service.Execute(CreateShipmentCommand{OrderID: "order-001"})
	if err != domain.ErrOrderNotShippable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotShippable, err)
	}
}
