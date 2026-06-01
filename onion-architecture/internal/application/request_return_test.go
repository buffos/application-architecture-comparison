package application

import (
	"testing"
	"time"

	"onion-architecture/internal/domain"
	timeinfra "onion-architecture/internal/infrastructure/services/time"
)

type stubReturnRequestStore struct {
	saved domain.ReturnRequest
	found domain.ReturnRequest
	list  []domain.ReturnRequest
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

func (s *stubReturnRequestStore) ListByStatus(status string) ([]domain.ReturnRequest, error) {
	result := make([]domain.ReturnRequest, 0)
	for _, request := range s.list {
		if request.Status == status {
			result = append(result, request)
		}
	}

	return result, nil
}

func TestRequestReturnServiceCreatesRequestedReturnForShippedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusShipped,
			ShippedAt:  time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
		},
	}
	returns := &stubReturnRequestStore{}
	clock := timeinfra.NewFixedClock(time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC))

	service := NewRequestReturnService(orders, returns, clock)

	result, err := service.Execute(RequestReturnCommand{
		OrderID:     "order-001",
		Reason:      "damaged on arrival",
		RequestedBy: "customer-001",
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

	if !returns.saved.RequestedAt.Equal(clock.Now()) {
		t.Fatalf("expected requested at %v, got %v", clock.Now(), returns.saved.RequestedAt)
	}

	if returns.saved.RequestedBy != "customer-001" {
		t.Fatalf("expected requested by customer-001, got %s", returns.saved.RequestedBy)
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
	clock := timeinfra.NewFixedClock(time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC))

	service := NewRequestReturnService(orders, returns, clock)

	_, err := service.Execute(RequestReturnCommand{
		OrderID:     "order-001",
		Reason:      "changed mind",
		RequestedBy: "customer-001",
	})
	if err != domain.ErrOrderNotReturnable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotReturnable, err)
	}
}

func TestRequestReturnServiceRejectsMissingActor(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:        "order-001",
			Status:    domain.OrderStatusShipped,
			ShippedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
		},
	}
	returns := &stubReturnRequestStore{}
	clock := timeinfra.NewFixedClock(time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC))

	service := NewRequestReturnService(orders, returns, clock)

	_, err := service.Execute(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged on arrival",
	})
	if err != domain.ErrActorRequired {
		t.Fatalf("expected %v, got %v", domain.ErrActorRequired, err)
	}
}
