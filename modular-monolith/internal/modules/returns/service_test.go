package returns

import (
	"testing"
	"time"

	"modular-monolith/internal/modules/idempotency"
	"modular-monolith/internal/modules/inventory"
	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/payments"
	"modular-monolith/internal/modules/returneligibility"
)

type stubRepository struct {
	saved ReturnRequest
}

func (r *stubRepository) Save(request ReturnRequest) error {
	r.saved = request
	return nil
}

func (r *stubRepository) FindByID(id string) (ReturnRequest, error) {
	return r.saved, nil
}

func (r *stubRepository) ListByStatus(status string) ([]ReturnRequest, error) {
	if status == "" || r.saved.Status == status {
		return []ReturnRequest{r.saved}, nil
	}

	return []ReturnRequest{}, nil
}

type stubOrderSource struct {
	order orders.ReturnableOrder
	err   error
}

func (s stubOrderSource) GetReturnableOrder(orderID string) (orders.ReturnableOrder, error) {
	if s.err != nil {
		return orders.ReturnableOrder{}, s.err
	}

	return s.order, nil
}

type stubRefunder struct {
	request payments.RefundRequest
	err     error
}

type stubRestocker struct {
	items []inventory.RestockItem
	err   error
}

type stubEligibility struct {
	allows bool
}

type stubClock struct {
	now time.Time
}

type stubIdempotencyStore struct {
	results map[string]idempotency.Result
}

func (s *stubRefunder) Refund(request payments.RefundRequest) error {
	if s.err != nil {
		return s.err
	}

	s.request = request
	return nil
}

func (s *stubRestocker) Restock(items []inventory.RestockItem) error {
	if s.err != nil {
		return s.err
	}

	s.items = append([]inventory.RestockItem(nil), items...)
	return nil
}

func (s stubEligibility) Allows(request returneligibility.ReviewRequest) bool {
	return s.allows
}

func (c stubClock) Now() time.Time {
	return c.now
}

func (s *stubIdempotencyStore) Find(key string) (idempotency.Result, bool, error) {
	result, ok := s.results[key]
	return result, ok, nil
}

func (s *stubIdempotencyStore) Save(key string, result idempotency.Result) error {
	s.results[key] = result
	return nil
}

func TestRequestReturnStoresRequestedReturn(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	clock := stubClock{now: time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC)}
	idempotencyStore := &stubIdempotencyStore{results: map[string]idempotency.Result{}}
	service := NewService(repository, stubOrderSource{
		order: orders.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
			Lines: []orders.ReturnableOrderLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000, ReturnWindowDays: 30},
			},
		},
	}, stubEligibility{allows: true}, restocker, idempotencyStore, refunder, clock)

	result, err := service.RequestReturn(RequestReturnCommand{
		OrderID:     "order-001",
		Reason:      "damaged item",
		RequestedBy: "customer-001",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != ReturnRequestStatusRequested {
		t.Fatalf("expected %s, got %s", ReturnRequestStatusRequested, result.Status)
	}

	if refunder.request.Amount != 0 {
		t.Fatalf("expected no refund during request, got %d", refunder.request.Amount)
	}

	if len(restocker.items) != 0 {
		t.Fatalf("expected no restock during request, got %+v", restocker.items)
	}

	if !repository.saved.RequestedAt.Equal(clock.now) {
		t.Fatalf("expected requested time to be recorded")
	}

	if repository.saved.RequestedBy != "customer-001" {
		t.Fatalf("expected requested by customer-001, got %s", repository.saved.RequestedBy)
	}
}

func TestRequestReturnRejectsNonReturnableOrder(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	idempotencyStore := &stubIdempotencyStore{results: map[string]idempotency.Result{}}
	service := NewService(repository, stubOrderSource{
		err: orders.ErrOrderNotReturnable,
	}, stubEligibility{allows: true}, restocker, idempotencyStore, refunder, stubClock{})

	_, err := service.RequestReturn(RequestReturnCommand{
		OrderID:     "order-001",
		Reason:      "damaged item",
		RequestedBy: "customer-001",
	})
	if err != orders.ErrOrderNotReturnable {
		t.Fatalf("expected %v, got %v", orders.ErrOrderNotReturnable, err)
	}
}

func TestRequestReturnRejectsMissingActor(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	clock := stubClock{now: time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC)}
	idempotencyStore := &stubIdempotencyStore{results: map[string]idempotency.Result{}}
	service := NewService(repository, stubOrderSource{
		order: orders.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
			Lines: []orders.ReturnableOrderLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000, ReturnWindowDays: 30},
			},
		},
	}, stubEligibility{allows: true}, restocker, idempotencyStore, refunder, clock)

	_, err := service.RequestReturn(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged item",
	})
	if err != ErrActorRequired {
		t.Fatalf("expected %v, got %v", ErrActorRequired, err)
	}
}

func TestRequestReturnStopsWhenRestockFails(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{err: inventory.ErrStockNotFound}
	idempotencyStore := &stubIdempotencyStore{results: map[string]idempotency.Result{}}
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: true}, restocker, idempotencyStore, refunder, stubClock{})
	repository.saved, _ = NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000, ReturnWindowDays: 30},
		},
	}, "damaged item", time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC), "customer-001")

	_, err := service.AcceptReturn(ReviewReturnCommand{
		ReturnRequestID: repository.saved.ID,
		IdempotencyKey:  "accept-1",
		ActorID:         "agent-001",
		ReviewNote:      "warehouse restock failed",
	})
	if err != inventory.ErrStockNotFound {
		t.Fatalf("expected %v, got %v", inventory.ErrStockNotFound, err)
	}
}

func TestAcceptReturnRefundsRestocksAndStoresUpdatedStatus(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	idempotencyStore := &stubIdempotencyStore{results: map[string]idempotency.Result{}}
	repository.saved, _ = NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000, ReturnWindowDays: 30},
		},
	}, "damaged item", time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC), "customer-001")
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: true}, restocker, idempotencyStore, refunder, stubClock{})

	result, err := service.AcceptReturn(ReviewReturnCommand{ReturnRequestID: repository.saved.ID, IdempotencyKey: "accept-1", ActorID: "agent-001", ReviewNote: "accepted after inspection"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != ReturnRequestStatusRefunded {
		t.Fatalf("expected %s, got %s", ReturnRequestStatusRefunded, result.Status)
	}

	if refunder.request.Amount != 30000 {
		t.Fatalf("expected refund amount 30000, got %d", refunder.request.Amount)
	}

	if len(restocker.items) != 1 || restocker.items[0].Quantity != 2 {
		t.Fatalf("expected restock quantity 2, got %+v", restocker.items)
	}

	if repository.saved.ReviewedBy != "agent-001" || repository.saved.ProcessedBy != "agent-001" {
		t.Fatalf("expected review and process actor to be recorded")
	}
}

func TestRejectReturnStoresRejectedStatus(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	idempotencyStore := &stubIdempotencyStore{results: map[string]idempotency.Result{}}
	repository.saved, _ = NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000, ReturnWindowDays: 30},
		},
	}, "damaged item", time.Date(2026, 6, 12, 12, 0, 0, 0, time.UTC), "customer-001")
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: true}, restocker, idempotencyStore, refunder, stubClock{})

	result, err := service.RejectReturn(ReviewReturnCommand{ReturnRequestID: repository.saved.ID, IdempotencyKey: "reject-1", ActorID: "agent-002", ReviewNote: "rejected on inspection"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected %s, got %s", ReturnRequestStatusRejected, result.Status)
	}

	if refunder.request.Amount != 0 {
		t.Fatalf("expected no refund on rejection, got %d", refunder.request.Amount)
	}

	if repository.saved.ReviewedBy != "agent-002" {
		t.Fatalf("expected reviewer to be recorded")
	}
}

func TestAcceptReturnRejectsWhenPolicyBlocksEligibility(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	idempotencyStore := &stubIdempotencyStore{results: map[string]idempotency.Result{}}
	repository.saved, _ = NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		ShippedAt:  time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC),
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000, ReturnWindowDays: 30},
		},
	}, "outside return window", time.Date(2026, 7, 12, 12, 0, 0, 0, time.UTC), "customer-001")
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: false}, restocker, idempotencyStore, refunder, stubClock{})

	result, err := service.AcceptReturn(ReviewReturnCommand{ReturnRequestID: repository.saved.ID, IdempotencyKey: "accept-2", ActorID: "agent-003", ReviewNote: "window exceeded"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected %s, got %s", ReturnRequestStatusRejected, result.Status)
	}

	if refunder.request.Amount != 0 {
		t.Fatalf("expected no refund when policy blocks return, got %d", refunder.request.Amount)
	}

	if len(restocker.items) != 0 {
		t.Fatalf("expected no restock when policy blocks return, got %+v", restocker.items)
	}

	if repository.saved.ReviewedBy != "agent-003" {
		t.Fatalf("expected policy reviewer to be recorded")
	}
}

func TestAcceptReturnReusesStoredIdempotentResult(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	idempotencyStore := &stubIdempotencyStore{
		results: map[string]idempotency.Result{
			"accept-1": {
				ReturnRequestID: "return-001",
				OrderID:         "order-001",
				CustomerID:      "customer-001",
				Status:          ReturnRequestStatusRefunded,
				LineCount:       1,
			},
		},
	}
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: true}, restocker, idempotencyStore, refunder, stubClock{})

	result, err := service.AcceptReturn(ReviewReturnCommand{ReturnRequestID: "return-001", IdempotencyKey: "accept-1", ActorID: "agent-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != ReturnRequestStatusRefunded {
		t.Fatalf("expected refunded result, got %s", result.Status)
	}

	if refunder.request.Amount != 0 || len(restocker.items) != 0 {
		t.Fatalf("expected no side effects on idempotent replay")
	}
}

func TestRejectReturnReusesStoredIdempotentResult(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	idempotencyStore := &stubIdempotencyStore{
		results: map[string]idempotency.Result{
			"reject-1": {
				ReturnRequestID: "return-001",
				OrderID:         "order-001",
				CustomerID:      "customer-001",
				Status:          ReturnRequestStatusRejected,
				LineCount:       1,
			},
		},
	}
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: true}, restocker, idempotencyStore, refunder, stubClock{})

	result, err := service.RejectReturn(ReviewReturnCommand{ReturnRequestID: "return-001", IdempotencyKey: "reject-1", ActorID: "agent-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected rejected result, got %s", result.Status)
	}

	if refunder.request.Amount != 0 || len(restocker.items) != 0 {
		t.Fatalf("expected no side effects on idempotent replay")
	}
}
