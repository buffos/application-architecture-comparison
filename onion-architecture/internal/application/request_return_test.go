package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubReturnRequestStore struct {
	saved domain.ReturnRequest
	found domain.ReturnRequest
}

func (s *stubReturnRequestStore) Save(request domain.ReturnRequest) error {
	s.saved = request
	return nil
}

func (s *stubReturnRequestStore) FindByID(id string) (domain.ReturnRequest, error) {
	if s.found.ID == "" {
		return domain.ReturnRequest{}, domain.ErrReturnRequestNotFound
	}

	return s.found, nil
}

func TestRequestReturnServiceCreatesRequestedReturnForShippedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusShipped,
		},
	}
	returns := &stubReturnRequestStore{}

	service := NewRequestReturnService(orders, returns)

	result, err := service.Execute(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged on arrival",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRequested {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRequested, result.Status)
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

	service := NewRequestReturnService(orders, returns)

	_, err := service.Execute(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "changed mind",
	})
	if err != domain.ErrOrderNotReturnable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotReturnable, err)
	}
}
