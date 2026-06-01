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

	service := NewRequestReturnService(orders, returns, stubRefundGateway{})

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

	service := NewRequestReturnService(orders, returns, stubRefundGateway{})

	_, err := service.Execute(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "changed mind",
	})
	if err != domain.ErrOrderNotReturnable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotReturnable, err)
	}
}
