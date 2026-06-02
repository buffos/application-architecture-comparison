package returns

import (
	"testing"
	"time"

	"modular-monolith/internal/modules/idempotency"
)

type stubQueryRepository struct {
	requests map[string]ReturnRequest
}

func (r *stubQueryRepository) Save(request ReturnRequest) error {
	if r.requests == nil {
		r.requests = make(map[string]ReturnRequest)
	}
	r.requests[request.ID] = request
	return nil
}

func (r *stubQueryRepository) FindByID(id string) (ReturnRequest, error) {
	request, ok := r.requests[id]
	if !ok {
		return ReturnRequest{}, ErrReturnRequestNotFound
	}
	return request, nil
}

func (r *stubQueryRepository) ListByStatus(status string) ([]ReturnRequest, error) {
	list := make([]ReturnRequest, 0, len(r.requests))
	for _, request := range r.requests {
		if status == "" || request.Status == status {
			list = append(list, request)
		}
	}
	return list, nil
}

func newQueryService(repository Repository) Service {
	return NewService(
		repository,
		stubOrderSource{},
		stubEligibility{allows: true},
		&stubRestocker{},
		&stubIdempotencyStore{results: map[string]idempotency.Result{}},
		&stubRefunder{},
		stubClock{},
	)
}

func TestGetReturnRequestLoadsStoredRequest(t *testing.T) {
	repository := &stubQueryRepository{requests: map[string]ReturnRequest{}}
	request, _ := NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, ShippedQuantity: 2, UnitPrice: 15000, ReturnWindowDays: 30},
		},
	}, nil, "damaged item", time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC), "customer-001")
	request.ReviewedBy = "agent-001"
	request.ProcessedBy = "agent-001"
	_ = repository.Save(request)
	service := newQueryService(repository)

	result, err := service.GetReturnRequest(GetReturnRequestQuery{ReturnRequestID: request.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ReturnRequestID != request.ID || result.RequestedBy != "customer-001" {
		t.Fatalf("expected stored request details to be returned")
	}
}

func TestListReturnRequestsFiltersByStatus(t *testing.T) {
	repository := &stubQueryRepository{requests: map[string]ReturnRequest{}}
	requested, _ := NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, ShippedQuantity: 2, UnitPrice: 15000, ReturnWindowDays: 30},
		},
	}, nil, "damaged item", time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC), "customer-001")
	refunded, _ := NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-002",
		CustomerID: "customer-002",
		ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-002", ProductName: "Chair", ProductCategory: "Standard", Quantity: 1, ShippedQuantity: 1, UnitPrice: 5000, ReturnWindowDays: 30},
		},
	}, nil, "damaged item", time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC), "customer-002")
	_ = refunded.Refund("agent-001", "agent-001", "accepted")
	_ = repository.Save(requested)
	_ = repository.Save(refunded)
	service := newQueryService(repository)

	result, err := service.ListReturnRequests(ListReturnRequestsQuery{Status: ReturnRequestStatusRequested})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].Status != ReturnRequestStatusRequested {
		t.Fatalf("expected one requested return, got %+v", result)
	}
}
