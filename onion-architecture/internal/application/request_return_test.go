package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubReturnRequestStore struct {
	saved domain.ReturnRequest
}

func (s *stubReturnRequestStore) Save(request domain.ReturnRequest) error {
	s.saved = request
	return nil
}

type stubRefundGateway struct {
	err error
}

func (g stubRefundGateway) Refund(order domain.Order) error {
	return g.err
}

type stubInventoryRestock struct {
	items []domain.InventoryRestockItem
	err   error
}

func (s *stubInventoryRestock) Restock(items []domain.InventoryRestockItem) error {
	if s.err != nil {
		return s.err
	}

	s.items = items
	return nil
}

func TestRequestReturnServiceCreatesRefundedReturnForShippedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusShipped,
		},
	}
	returns := &stubReturnRequestStore{}
	restock := &stubInventoryRestock{}

	service := NewRequestReturnService(orders, returns, stubRefundGateway{}, restock)

	result, err := service.Execute(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged on arrival",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRefunded {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRefunded, result.Status)
	}

	if returns.saved.OrderID != "order-001" {
		t.Fatalf("expected saved order id order-001, got %s", returns.saved.OrderID)
	}

	if len(restock.items) != 0 {
		t.Fatalf("expected no restock items for empty order lines, got %d", len(restock.items))
	}
}

func TestRequestReturnServiceRejectsNonShippedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPaid,
		},
	}
	returns := &stubReturnRequestStore{}
	restock := &stubInventoryRestock{}

	service := NewRequestReturnService(orders, returns, stubRefundGateway{}, restock)

	_, err := service.Execute(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "changed mind",
	})
	if err != domain.ErrOrderNotReturnable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotReturnable, err)
	}
}

func TestRequestReturnServiceRestocksInventoryFromReturnedLines(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusShipped,
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
	returns := &stubReturnRequestStore{}
	restock := &stubInventoryRestock{}

	service := NewRequestReturnService(orders, returns, stubRefundGateway{}, restock)

	_, err := service.Execute(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged on arrival",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected one restock item, got %d", len(restock.items))
	}

	if restock.items[0].Quantity != 2 {
		t.Fatalf("expected restock quantity 2, got %d", restock.items[0].Quantity)
	}
}
