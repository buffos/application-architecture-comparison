package returns

import (
	"errors"
	"testing"
	"time"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	saved   ReturnRequest
	saveErr error
}

func (r *stubRepository) FindByID(id string) (ReturnRequest, error) {
	if r.saved.ID == id {
		return r.saved, nil
	}

	return ReturnRequest{}, ErrReturnRequestNotFound
}

func (r *stubRepository) ListByStatus(status string) ([]ReturnRequest, error) {
	if r.saved.ID == "" {
		return []ReturnRequest{}, nil
	}

	if status == "" || r.saved.Status == status {
		return []ReturnRequest{r.saved}, nil
	}

	return []ReturnRequest{}, nil
}

func (r *stubRepository) Save(request ReturnRequest) error {
	if r.saveErr != nil {
		return r.saveErr
	}

	r.saved = request
	return nil
}

type stubReturnableOrderProvider struct {
	order kernel.ReturnableOrder
	err   error
}

func (p stubReturnableOrderProvider) GetReturnableOrder(orderID string) (kernel.ReturnableOrder, error) {
	return p.order, p.err
}

type stubEligibilityPolicy struct {
	allowed bool
}

func (p stubEligibilityPolicy) Allows(review kernel.ReturnEligibilityReview) bool {
	return p.allowed
}

type stubPaymentRefund struct {
	orderID string
	amount  int
	err     error
}

func (p *stubPaymentRefund) Refund(orderID string, amount int) error {
	if p.err != nil {
		return p.err
	}

	p.orderID = orderID
	p.amount = amount
	return nil
}

type stubInventoryRestock struct {
	items []kernel.InventoryReservationItem
	err   error
}

func (s *stubInventoryRestock) Restock(items []kernel.InventoryReservationItem) error {
	if s.err != nil {
		return s.err
	}

	s.items = append([]kernel.InventoryReservationItem(nil), items...)
	return nil
}

type stubClock struct {
	now time.Time
}

func (c stubClock) Now() time.Time {
	return c.now
}

type stubIdempotencyStore struct {
	results map[string]kernel.IdempotencyResult
	saveErr error
}

func (s *stubIdempotencyStore) Find(key string) (kernel.IdempotencyResult, bool, error) {
	result, ok := s.results[key]
	return result, ok, nil
}

func (s *stubIdempotencyStore) Save(key string, result kernel.IdempotencyResult) error {
	if s.saveErr != nil {
		return s.saveErr
	}

	if s.results == nil {
		s.results = map[string]kernel.IdempotencyResult{}
	}

	s.results[key] = result
	return nil
}

func TestRequestReturnStoresRequestedReturnWithoutRefundOrRestock(t *testing.T) {
	repository := &stubRepository{}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}
	service := NewService(repository, stubReturnableOrderProvider{
		order: kernel.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			ShippedAt:  time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			Lines: []kernel.ReturnableOrderLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}, stubClock{now: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)}, stubEligibilityPolicy{allowed: true}, idempotency, refunds, restock)

	result, err := service.RequestReturn(kernel.RequestReturnCommand{
		OrderID:     "order-001",
		Reason:      "damaged item",
		RequestedBy: "customer-001",
	})
	if err != nil {
		t.Fatalf("expected request return to succeed, got %v", err)
	}

	if result.Status != ReturnRequestStatusRequested {
		t.Fatalf("expected requested status, got %s", result.Status)
	}

	if repository.saved.Status != ReturnRequestStatusRequested {
		t.Fatalf("expected saved request to be requested, got %s", repository.saved.Status)
	}

	if repository.saved.RequestedAt.IsZero() {
		t.Fatalf("expected requested time to be recorded")
	}

	if repository.saved.RequestedBy != "customer-001" {
		t.Fatalf("expected requester to be recorded, got %s", repository.saved.RequestedBy)
	}

	if refunds.orderID != "" || refunds.amount != 0 {
		t.Fatalf("expected no refund during request, got %s %d", refunds.orderID, refunds.amount)
	}

	if len(restock.items) != 0 {
		t.Fatalf("expected no restock during request, got %d items", len(restock.items))
	}
}

func TestAcceptReturnRefundsRestocksAndUpdatesStatus(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			CustomerID:  "customer-001",
			ShippedAt:   time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			RequestedAt: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
			Status:      ReturnRequestStatusRequested,
			Lines: []ReturnLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, idempotency, refunds, restock)

	result, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "accept-001",
		ReviewedBy:      "agent-001",
		ProcessedBy:     "ops-001",
		ReviewNote:      "approved after inspection",
	})
	if err != nil {
		t.Fatalf("expected accept return to succeed, got %v", err)
	}

	if result.Status != ReturnRequestStatusRefunded {
		t.Fatalf("expected refunded status, got %s", result.Status)
	}

	if refunds.orderID != "order-001" || refunds.amount != 45000 {
		t.Fatalf("expected refund for order-001 amount 45000, got %s %d", refunds.orderID, refunds.amount)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected 1 restock item, got %d", len(restock.items))
	}

	if repository.saved.Status != ReturnRequestStatusRefunded {
		t.Fatalf("expected saved request to be refunded, got %s", repository.saved.Status)
	}

	if repository.saved.ReviewedBy != "agent-001" || repository.saved.ProcessedBy != "ops-001" {
		t.Fatalf("expected review and processor metadata to be recorded, got %s %s", repository.saved.ReviewedBy, repository.saved.ProcessedBy)
	}

	if _, ok := idempotency.results["accept-001"]; !ok {
		t.Fatalf("expected idempotent result to be stored")
	}
}

func TestAcceptReturnStopsWhenRestockFails(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			CustomerID:  "customer-001",
			ShippedAt:   time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			RequestedAt: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
			Status:      ReturnRequestStatusRequested,
			Lines: []ReturnLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{err: errors.New("restock failed")}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}, refunds, restock)

	_, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "accept-002",
		ReviewedBy:      "agent-001",
		ProcessedBy:     "ops-001",
	})
	if err == nil || err.Error() != "restock failed" {
		t.Fatalf("expected restock failure, got %v", err)
	}

	if repository.saved.Status != ReturnRequestStatusRequested {
		t.Fatalf("expected request to remain requested, got %s", repository.saved.Status)
	}
}

func TestRejectReturnUpdatesStatusWithoutRefundOrRestock(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			CustomerID:  "customer-001",
			ShippedAt:   time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			RequestedAt: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
			Status:      ReturnRequestStatusRequested,
			Lines: []ReturnLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, idempotency, refunds, restock)

	result, err := service.RejectReturn(kernel.RejectReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "reject-001",
		ReviewedBy:      "agent-002",
		ReviewNote:      "missing evidence",
	})
	if err != nil {
		t.Fatalf("expected reject return to succeed, got %v", err)
	}

	if result.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected rejected status, got %s", result.Status)
	}

	if repository.saved.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected saved request to be rejected, got %s", repository.saved.Status)
	}

	if repository.saved.ReviewedBy != "agent-002" || repository.saved.ProcessedBy != "" {
		t.Fatalf("expected reviewer metadata without processor, got %s %s", repository.saved.ReviewedBy, repository.saved.ProcessedBy)
	}

	if refunds.orderID != "" || refunds.amount != 0 {
		t.Fatalf("expected no refund during rejection, got %s %d", refunds.orderID, refunds.amount)
	}

	if len(restock.items) != 0 {
		t.Fatalf("expected no restock during rejection, got %d items", len(restock.items))
	}

	if _, ok := idempotency.results["reject-001"]; !ok {
		t.Fatalf("expected reject result to be stored idempotently")
	}
}

func TestAcceptReturnRejectsWhenPolicyBlocks(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			CustomerID:  "customer-001",
			Reason:      "damaged item",
			ShippedAt:   time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			RequestedAt: time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC),
			Status:      ReturnRequestStatusRequested,
			Lines: []ReturnLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: false}, idempotency, refunds, restock)

	result, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "accept-003",
		ReviewedBy:      "agent-003",
		ReviewNote:      "window expired",
	})
	if err != nil {
		t.Fatalf("expected policy-blocked accept to succeed with rejection, got %v", err)
	}

	if result.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected rejected status, got %s", result.Status)
	}

	if repository.saved.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected saved request to be rejected, got %s", repository.saved.Status)
	}

	if repository.saved.ReviewedBy != "agent-003" || repository.saved.ProcessedBy != "" {
		t.Fatalf("expected reviewer-only metadata on blocked return, got %s %s", repository.saved.ReviewedBy, repository.saved.ProcessedBy)
	}

	if refunds.orderID != "" || refunds.amount != 0 {
		t.Fatalf("expected no refund when policy blocks, got %s %d", refunds.orderID, refunds.amount)
	}

	if len(restock.items) != 0 {
		t.Fatalf("expected no restock when policy blocks, got %d items", len(restock.items))
	}
}

func TestRequestReturnRejectsMissingRequester(t *testing.T) {
	repository := &stubRepository{}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}
	service := NewService(repository, stubReturnableOrderProvider{
		order: kernel.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			ShippedAt:  time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			Lines: []kernel.ReturnableOrderLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}, stubClock{now: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)}, stubEligibilityPolicy{allowed: true}, idempotency, refunds, restock)

	_, err := service.RequestReturn(kernel.RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged item",
	})
	if err != ErrActorRequired {
		t.Fatalf("expected missing requester to be rejected, got %v", err)
	}
}

func TestAcceptReturnRejectsMissingActors(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			CustomerID:  "customer-001",
			ShippedAt:   time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			RequestedAt: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
			Status:      ReturnRequestStatusRequested,
			Lines: []ReturnLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}, &stubPaymentRefund{}, &stubInventoryRestock{})

	_, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "accept-004",
		ReviewedBy:      "agent-001",
	})
	if err != ErrActorRequired {
		t.Fatalf("expected missing processor to be rejected, got %v", err)
	}
}

func TestRejectReturnRejectsMissingReviewer(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			CustomerID:  "customer-001",
			ShippedAt:   time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			RequestedAt: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
			Status:      ReturnRequestStatusRequested,
			Lines: []ReturnLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}, &stubPaymentRefund{}, &stubInventoryRestock{})

	_, err := service.RejectReturn(kernel.RejectReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "reject-002",
	})
	if err != ErrActorRequired {
		t.Fatalf("expected missing reviewer to be rejected, got %v", err)
	}
}

func TestAcceptReturnReusesStoredResultOnRetry(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{ID: "return-001"},
	}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{
		results: map[string]kernel.IdempotencyResult{
			"accept-005": {
				ReturnRequestID: "return-001",
				OrderID:         "order-001",
				CustomerID:      "customer-001",
				Status:          ReturnRequestStatusRefunded,
				LineCount:       1,
			},
		},
	}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, idempotency, refunds, restock)

	result, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "accept-005",
		ReviewedBy:      "agent-001",
		ProcessedBy:     "ops-001",
	})
	if err != nil {
		t.Fatalf("expected idempotent retry to succeed, got %v", err)
	}

	if result.Status != ReturnRequestStatusRefunded {
		t.Fatalf("expected refunded status, got %s", result.Status)
	}

	if refunds.orderID != "" || len(restock.items) != 0 {
		t.Fatalf("expected no side effects on replay")
	}
}

func TestRejectReturnReusesStoredResultOnRetry(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{ID: "return-001"},
	}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	idempotency := &stubIdempotencyStore{
		results: map[string]kernel.IdempotencyResult{
			"reject-003": {
				ReturnRequestID: "return-001",
				OrderID:         "order-001",
				CustomerID:      "customer-001",
				Status:          ReturnRequestStatusRejected,
				LineCount:       1,
			},
		},
	}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, idempotency, refunds, restock)

	result, err := service.RejectReturn(kernel.RejectReturnCommand{
		ReturnRequestID: "return-001",
		IdempotencyKey:  "reject-003",
		ReviewedBy:      "agent-002",
	})
	if err != nil {
		t.Fatalf("expected idempotent retry to succeed, got %v", err)
	}

	if result.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected rejected status, got %s", result.Status)
	}

	if refunds.orderID != "" || len(restock.items) != 0 {
		t.Fatalf("expected no side effects on reject replay")
	}
}

func TestAcceptReturnRejectsMissingIdempotencyKey(t *testing.T) {
	service := NewService(&stubRepository{}, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}, &stubPaymentRefund{}, &stubInventoryRestock{})

	_, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
		ReviewedBy:      "agent-001",
		ProcessedBy:     "ops-001",
	})
	if err != kernel.ErrIdempotencyKeyRequired {
		t.Fatalf("expected missing idempotency key error, got %v", err)
	}
}

func TestRejectReturnRejectsMissingIdempotencyKey(t *testing.T) {
	service := NewService(&stubRepository{}, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}, &stubPaymentRefund{}, &stubInventoryRestock{})

	_, err := service.RejectReturn(kernel.RejectReturnCommand{
		ReturnRequestID: "return-001",
		ReviewedBy:      "agent-002",
	})
	if err != kernel.ErrIdempotencyKeyRequired {
		t.Fatalf("expected missing idempotency key error, got %v", err)
	}
}

func TestGetReturnRequest(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			CustomerID:  "customer-001",
			Reason:      "damaged item",
			RequestedBy: "customer-001",
			ReviewedBy:  "agent-001",
			ProcessedBy: "ops-001",
			ReviewNote:  "approved after inspection",
			Status:      ReturnRequestStatusRefunded,
			Lines: []ReturnLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}, &stubPaymentRefund{}, &stubInventoryRestock{})

	result, err := service.GetReturnRequest(kernel.GetReturnRequestQuery{
		ReturnRequestID: "return-001",
	})
	if err != nil {
		t.Fatalf("expected get return request to succeed, got %v", err)
	}

	if result.ReturnRequestID != "return-001" || result.RequestedBy != "customer-001" {
		t.Fatalf("unexpected return details %+v", result)
	}
}

func TestListReturnRequestsByStatus(t *testing.T) {
	repository := &stubRepository{
		saved: ReturnRequest{
			ID:      "return-001",
			OrderID: "order-001",
			Status:  ReturnRequestStatusRequested,
			Lines: []ReturnLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, &stubIdempotencyStore{results: map[string]kernel.IdempotencyResult{}}, &stubPaymentRefund{}, &stubInventoryRestock{})

	result, err := service.ListReturnRequests(kernel.ListReturnRequestsQuery{
		Status: ReturnRequestStatusRequested,
	})
	if err != nil {
		t.Fatalf("expected list return requests to succeed, got %v", err)
	}

	if len(result) != 1 || result[0].ReturnRequestID != "return-001" {
		t.Fatalf("unexpected list result %+v", result)
	}
}
