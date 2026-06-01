package application

import (
	"testing"
	"time"

	"onion-architecture/internal/domain"
	returneligibility "onion-architecture/internal/infrastructure/policies/returneligibility"
)

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

func TestAcceptReturnServiceRefundsAndRestocksAcceptedReturn(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			Status:     domain.OrderStatusShipped,
			ShippedAt:  time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:       "sku-002",
					Quantity:         2,
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
			Reason:      "damaged on arrival",
			RequestedAt: time.Date(2026, 6, 5, 10, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
		},
	}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{entries: make(map[string]string)}

	service := NewAcceptReturnService(orders, returns, stubReturnEligibilityPolicy{eligible: true}, idempotency, stubRefundGateway{}, restock)

	result, err := service.Execute(AcceptReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "accept-001",
		ReviewedBy:      "agent-007",
		ReviewNote:      "confirmed damage",
		ProcessedBy:     "finance-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRefunded {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRefunded, result.Status)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected one restock item, got %d", len(restock.items))
	}

	if returns.saved.ReviewedBy != "agent-007" {
		t.Fatalf("expected reviewed by agent-007, got %s", returns.saved.ReviewedBy)
	}

	if returns.saved.ProcessedBy != "finance-001" {
		t.Fatalf("expected processed by finance-001, got %s", returns.saved.ProcessedBy)
	}
}

func TestRejectReturnServiceRejectsRequestedReturn(t *testing.T) {
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			Status:      domain.ReturnRequestStatusRequested,
			RequestedBy: "customer-001",
		},
	}

	idempotency := &stubIdempotencyStore{entries: make(map[string]string)}

	service := NewRejectReturnService(returns, idempotency)

	result, err := service.Execute(RejectReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "reject-001",
		ReviewedBy:      "agent-008",
		ReviewNote:      "opened and used",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRejected {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRejected, result.Status)
	}

	if returns.saved.ReviewedBy != "agent-008" {
		t.Fatalf("expected reviewed by agent-008, got %s", returns.saved.ReviewedBy)
	}
}

type stubReturnEligibilityPolicy struct {
	eligible bool
	err      error
}

func (p stubReturnEligibilityPolicy) IsEligible(request domain.ReturnRequest, order domain.Order) (bool, error) {
	if p.err != nil {
		return false, p.err
	}

	return p.eligible, nil
}

func TestAcceptReturnServiceLeavesRequestUnchangedWhenPolicyBlocksIt(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:     "order-001",
			Status: domain.OrderStatusShipped,
			ShippedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:       "sku-002",
					Quantity:         2,
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
			Reason:      "outside return window",
			RequestedAt: time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
		},
	}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{entries: make(map[string]string)}

	service := NewAcceptReturnService(orders, returns, stubReturnEligibilityPolicy{eligible: false}, idempotency, stubRefundGateway{}, restock)

	result, err := service.Execute(AcceptReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "accept-002",
		ReviewedBy:      "agent-009",
		ReviewNote:      "policy blocked",
		ProcessedBy:     "finance-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRequested {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRequested, result.Status)
	}

	if len(restock.items) != 0 {
		t.Fatalf("expected no restock items when policy blocks return, got %d", len(restock.items))
	}
}

func TestAcceptReturnServiceAppliesRealWindowPolicy(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:        "order-002",
			Status:    domain.OrderStatusShipped,
			ShippedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:       "sku-002",
					Quantity:         1,
					ReturnWindowDays: 30,
				},
			},
		},
	}
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:          "return-002",
			OrderID:     "order-002",
			Status:      domain.ReturnRequestStatusRequested,
			Reason:      "damaged on arrival",
			RequestedAt: time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
		},
	}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{entries: make(map[string]string)}

	service := NewAcceptReturnService(orders, returns, returneligibility.NewWindowPolicy(), idempotency, stubRefundGateway{}, restock)

	result, err := service.Execute(AcceptReturnCommand{
		ReturnRequestID: "return-002",
		IdempotencyKey:  "accept-003",
		ReviewedBy:      "agent-010",
		ReviewNote:      "within window",
		ProcessedBy:     "finance-002",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRefunded {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRefunded, result.Status)
	}
}

func TestAcceptReturnServiceBlocksOutOfWindowReturn(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:        "order-003",
			Status:    domain.OrderStatusShipped,
			ShippedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:       "sku-002",
					Quantity:         1,
					ReturnWindowDays: 30,
				},
			},
		},
	}
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:          "return-003",
			OrderID:     "order-003",
			Status:      domain.ReturnRequestStatusRequested,
			Reason:      "damaged on arrival",
			RequestedAt: time.Date(2026, 7, 5, 10, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
		},
	}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{entries: make(map[string]string)}

	service := NewAcceptReturnService(orders, returns, returneligibility.NewWindowPolicy(), idempotency, stubRefundGateway{}, restock)

	result, err := service.Execute(AcceptReturnCommand{
		ReturnRequestID: "return-003",
		IdempotencyKey:  "accept-004",
		ReviewedBy:      "agent-011",
		ReviewNote:      "outside window",
		ProcessedBy:     "finance-003",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRequested {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRequested, result.Status)
	}
}

type stubIdempotencyStore struct {
	entries map[string]string
}

func (s *stubIdempotencyStore) Get(key string) (string, bool, error) {
	status, ok := s.entries[key]
	return status, ok, nil
}

func (s *stubIdempotencyStore) Save(key string, status string) error {
	s.entries[key] = status
	return nil
}

func TestAcceptReturnServiceIsIdempotent(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:        "order-004",
			Status:    domain.OrderStatusShipped,
			ShippedAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC),
			Lines: []domain.OrderLine{
				{
					ProductSKU:       "sku-002",
					Quantity:         1,
					ReturnWindowDays: 30,
				},
			},
		},
	}
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:          "return-004",
			OrderID:     "order-004",
			Status:      domain.ReturnRequestStatusRequested,
			Reason:      "damaged on arrival",
			RequestedAt: time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
		},
	}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{entries: make(map[string]string)}

	service := NewAcceptReturnService(orders, returns, returneligibility.NewWindowPolicy(), idempotency, stubRefundGateway{}, restock)

	first, err := service.Execute(AcceptReturnCommand{
		ReturnRequestID: "return-004",
		IdempotencyKey:  "accept-005",
		ReviewedBy:      "agent-012",
		ReviewNote:      "within window",
		ProcessedBy:     "finance-004",
	})
	if err != nil {
		t.Fatalf("expected no error on first call, got %v", err)
	}

	second, err := service.Execute(AcceptReturnCommand{
		ReturnRequestID: "return-004",
		IdempotencyKey:  "accept-005",
		ReviewedBy:      "agent-012",
		ReviewNote:      "within window",
		ProcessedBy:     "finance-004",
	})
	if err != nil {
		t.Fatalf("expected no error on second call, got %v", err)
	}

	if first.Status != second.Status {
		t.Fatalf("expected same status on retry, got %s and %s", first.Status, second.Status)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected one restock execution, got %d", len(restock.items))
	}
}

func TestRejectReturnServiceIsIdempotent(t *testing.T) {
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:          "return-005",
			OrderID:     "order-005",
			Status:      domain.ReturnRequestStatusRequested,
			RequestedBy: "customer-001",
		},
	}
	idempotency := &stubIdempotencyStore{entries: make(map[string]string)}

	service := NewRejectReturnService(returns, idempotency)

	first, err := service.Execute(RejectReturnCommand{
		ReturnRequestID: "return-005",
		IdempotencyKey:  "reject-002",
		ReviewedBy:      "agent-013",
		ReviewNote:      "opened and used",
	})
	if err != nil {
		t.Fatalf("expected no error on first call, got %v", err)
	}

	second, err := service.Execute(RejectReturnCommand{
		ReturnRequestID: "return-005",
		IdempotencyKey:  "reject-002",
		ReviewedBy:      "agent-013",
		ReviewNote:      "opened and used",
	})
	if err != nil {
		t.Fatalf("expected no error on second call, got %v", err)
	}

	if first.Status != second.Status {
		t.Fatalf("expected same status on retry, got %s and %s", first.Status, second.Status)
	}
}
