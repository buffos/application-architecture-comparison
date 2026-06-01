package application

import (
	"testing"
	"time"

	"onion-architecture/internal/domain"
)

func TestRequestReturnServiceStoresOnlyRequestedReturnLines(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:        "order-001",
			Status:    domain.OrderStatusPartiallyShipped,
			ShippedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:       "sku-001",
					ProductCategory:  "Standard",
					Quantity:         3,
					ShippedQuantity:  2,
					ReturnWindowDays: 30,
				},
				{
					ProductSKU:       "sku-002",
					ProductCategory:  "CustomBuild",
					Quantity:         1,
					ShippedQuantity:  1,
					ReturnWindowDays: 30,
				},
			},
		},
	}
	returns := &stubReturnRequestStore{}
	clock := fixedClockForPartialReturn()
	service := NewRequestReturnService(orders, returns, clock)

	_, err := service.Execute(RequestReturnCommand{
		OrderID:     "order-001",
		Reason:      "damaged on arrival",
		RequestedBy: "customer-001",
		Lines: []RequestReturnLine{
			{ProductSKU: "sku-001", Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(returns.saved.Lines) != 1 {
		t.Fatalf("expected 1 return line, got %d", len(returns.saved.Lines))
	}

	if returns.saved.Lines[0].ProductSKU != "sku-001" || returns.saved.Lines[0].Quantity != 1 {
		t.Fatalf("unexpected return lines: %+v", returns.saved.Lines)
	}
}

func TestAcceptReturnServiceUpdatesOrderReturnedQuantityForPartialReturn(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:        "order-001",
			Status:    domain.OrderStatusPartiallyShipped,
			ShippedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:       "sku-001",
					ProductCategory:  "Standard",
					Quantity:         3,
					ShippedQuantity:  2,
					ReturnWindowDays: 30,
				},
				{
					ProductSKU:       "sku-002",
					ProductCategory:  "CustomBuild",
					Quantity:         1,
					ShippedQuantity:  1,
					ReturnWindowDays: 30,
				},
			},
		},
	}
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			Status:      domain.ReturnRequestStatusRequested,
			RequestedAt: time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
			Lines: []domain.ReturnRequestLine{
				{
					ProductSKU:       "sku-001",
					ProductCategory:  "Standard",
					Quantity:         1,
					ReturnWindowDays: 30,
				},
			},
		},
	}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{entries: make(map[string]string)}
	service := NewAcceptReturnService(orders, returns, stubReturnEligibilityPolicy{eligible: true}, idempotency, stubRefundGateway{}, restock)

	_, err := service.Execute(AcceptReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "accept-partial-001",
		ReviewedBy:      "agent-001",
		ReviewNote:      "partial accepted",
		ProcessedBy:     "finance-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(restock.items) != 1 || restock.items[0].Quantity != 1 {
		t.Fatalf("expected one restock item of quantity 1, got %+v", restock.items)
	}

	if orders.saved.Lines[0].ReturnedQuantity != 1 {
		t.Fatalf("expected returned quantity 1, got %d", orders.saved.Lines[0].ReturnedQuantity)
	}
}

func TestRequestReturnServiceRejectsReturnQuantityBeyondShippedMinusReturned(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:        "order-001",
			Status:    domain.OrderStatusPartiallyShipped,
			ShippedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:       "sku-001",
					ProductCategory:  "Standard",
					Quantity:         3,
					ShippedQuantity:  2,
					ReturnedQuantity: 1,
					ReturnWindowDays: 30,
				},
			},
		},
	}
	returns := &stubReturnRequestStore{}
	service := NewRequestReturnService(orders, returns, fixedClockForPartialReturn())

	_, err := service.Execute(RequestReturnCommand{
		OrderID:     "order-001",
		Reason:      "changed mind",
		RequestedBy: "customer-001",
		Lines: []RequestReturnLine{
			{ProductSKU: "sku-001", Quantity: 2},
		},
	})
	if err != domain.ErrReturnQuantityExceedsReturnable {
		t.Fatalf("expected %v, got %v", domain.ErrReturnQuantityExceedsReturnable, err)
	}
}

func fixedClockForPartialReturn() fixedClock {
	return fixedClock{now: time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC)}
}

type fixedClock struct {
	now time.Time
}

func (c fixedClock) Now() time.Time {
	return c.now
}
