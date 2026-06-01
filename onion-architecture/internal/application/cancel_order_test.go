package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubInventoryRelease struct {
	items []domain.InventoryReleaseItem
	err   error
}

func (s *stubInventoryRelease) Release(items []domain.InventoryReleaseItem) error {
	if s.err != nil {
		return s.err
	}

	s.items = items
	return nil
}

func TestCancelOrderServiceCancelsUnshippedOrderAndReleasesStock(t *testing.T) {
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
	inventory := &stubInventoryRelease{}

	service := NewCancelOrderService(orders, inventory)

	result, err := service.Execute(CancelOrderCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.OrderStatusCancelled {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusCancelled, result.Status)
	}

	if orders.saved.Status != domain.OrderStatusCancelled {
		t.Fatalf("expected saved status %s, got %s", domain.OrderStatusCancelled, orders.saved.Status)
	}

	if len(inventory.items) != 1 {
		t.Fatalf("expected one release item, got %d", len(inventory.items))
	}
}

func TestCancelOrderServiceRejectsShippedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusShipped,
		},
	}
	inventory := &stubInventoryRelease{}

	service := NewCancelOrderService(orders, inventory)

	_, err := service.Execute(CancelOrderCommand{OrderID: "order-001"})
	if err != domain.ErrOrderNotCancellable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotCancellable, err)
	}
}

func TestCancelOrderServiceRejectsPartiallyShippedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPartiallyShipped,
		},
	}
	inventory := &stubInventoryRelease{}

	service := NewCancelOrderService(orders, inventory)

	_, err := service.Execute(CancelOrderCommand{OrderID: "order-001"})
	if err != domain.ErrOrderNotCancellable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotCancellable, err)
	}
}
