package usecases

import (
	"testing"
	"time"

	"clean-architecture/internal/entities"
)

type stubReturnRequestEditor struct {
	request entities.ReturnRequest
	err     error
	saved   entities.ReturnRequest
}

func (g *stubReturnRequestEditor) FindByID(id string) (entities.ReturnRequest, error) {
	if g.err != nil {
		return entities.ReturnRequest{}, g.err
	}

	return g.request, nil
}

func (g *stubReturnRequestEditor) Save(request entities.ReturnRequest) error {
	g.saved = request
	return nil
}

type stubAcceptReturnOutput struct {
	output AcceptReturnOutput
}

func (o *stubAcceptReturnOutput) Present(output AcceptReturnOutput) error {
	o.output = output
	return nil
}

type stubReturnEligibilityPolicy struct {
	allowed bool
	err     error
}

func (p stubReturnEligibilityPolicy) CanAccept(order entities.Order, request entities.ReturnRequest) (bool, error) {
	if p.err != nil {
		return false, p.err
	}

	return p.allowed, nil
}

type stubRejectReturnOutput struct {
	output RejectReturnOutput
}

func (o *stubRejectReturnOutput) Present(output RejectReturnOutput) error {
	o.output = output
	return nil
}

func TestRequestReturnInteractorCreatesRequestedReturn(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusShipped,
			ShippedAt:     timePtr(time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)),
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2, ReturnWindowDays: 30},
			},
		},
	}
	returns := &stubReturnRequestWriter{}
	output := &stubRequestReturnOutput{}
	clock := stubClock{now: time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)}

	interactor := NewRequestReturnInteractor(orders, returns, clock, output)

	err := interactor.Execute(RequestReturnInput{OrderID: "order-001", Reason: "damaged item", RequestedBy: "customer-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.saved.Status != entities.ReturnRequestStatusRequested {
		t.Fatalf("expected status %s, got %s", entities.ReturnRequestStatusRequested, returns.saved.Status)
	}

	if returns.saved.RequestedBy != "customer-001" {
		t.Fatalf("expected requester customer-001, got %s", returns.saved.RequestedBy)
	}
}

func TestAcceptReturnInteractorRefundsAndRestocks(t *testing.T) {
	returns := &stubReturnRequestEditor{
		request: entities.ReturnRequest{
			ID:          "return-001",
			OrderID:     "order-001",
			Reason:      "damaged item",
			Status:      entities.ReturnRequestStatusRequested,
			RequestedAt: time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
		},
	}
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusShipped,
			ShippedAt:     timePtr(time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)),
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2, ReturnWindowDays: 30},
			},
		},
	}
	restock := &stubInventoryRestock{}
	output := &stubAcceptReturnOutput{}
	refunds := &stubRefundGateway{}
	idempotency := &stubIdempotencyStore{}

	interactor := NewAcceptReturnInteractor(idempotency, orders, returns, stubReturnEligibilityPolicy{allowed: true}, refunds, restock, output)

	err := interactor.Execute(AcceptReturnInput{ReturnRequestID: "return-001", IdempotencyKey: "accept-001", ReviewedBy: "reviewer-001", ProcessedBy: "finance-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.saved.Status != entities.ReturnRequestStatusRefunded {
		t.Fatalf("expected status %s, got %s", entities.ReturnRequestStatusRefunded, returns.saved.Status)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected 1 restock item, got %d", len(restock.items))
	}

	if returns.saved.ReviewedBy != "reviewer-001" {
		t.Fatalf("expected reviewer reviewer-001, got %s", returns.saved.ReviewedBy)
	}

	if returns.saved.ProcessedBy != "finance-001" {
		t.Fatalf("expected processor finance-001, got %s", returns.saved.ProcessedBy)
	}

	if refunds.calls != 1 {
		t.Fatalf("expected 1 refund call, got %d", refunds.calls)
	}

	if idempotency.records["accept-return:accept-001"] != "return-001" {
		t.Fatal("expected idempotency record to be saved for accepted return")
	}
}

func TestAcceptReturnInteractorBlocksPolicyRejectedReturn(t *testing.T) {
	returns := &stubReturnRequestEditor{
		request: entities.ReturnRequest{
			ID:          "return-003",
			OrderID:     "order-001",
			Reason:      "changed mind",
			Status:      entities.ReturnRequestStatusRequested,
			RequestedAt: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC),
			RequestedBy: "customer-001",
		},
	}
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusShipped,
			ShippedAt:     timePtr(time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)),
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2, ReturnWindowDays: 30},
			},
		},
	}
	restock := &stubInventoryRestock{}
	output := &stubAcceptReturnOutput{}
	idempotency := &stubIdempotencyStore{}

	interactor := NewAcceptReturnInteractor(idempotency, orders, returns, stubReturnEligibilityPolicy{allowed: false}, &stubRefundGateway{}, restock, output)

	err := interactor.Execute(AcceptReturnInput{ReturnRequestID: "return-003", IdempotencyKey: "accept-002", ReviewedBy: "reviewer-001", ProcessedBy: "finance-001"})
	if err != entities.ErrQuoteCannotTransition {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotTransition, err)
	}

	if returns.saved.ID != "" {
		t.Fatal("expected no saved return update when policy blocks acceptance")
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func TestRejectReturnInteractorPreventsRefundAndRestock(t *testing.T) {
	returns := &stubReturnRequestEditor{
		request: entities.ReturnRequest{
			ID:          "return-002",
			OrderID:     "order-001",
			Reason:      "changed mind",
			Status:      entities.ReturnRequestStatusRequested,
			RequestedBy: "customer-001",
		},
	}
	output := &stubRejectReturnOutput{}
	idempotency := &stubIdempotencyStore{}

	interactor := NewRejectReturnInteractor(idempotency, returns, output)

	err := interactor.Execute(RejectReturnInput{ReturnRequestID: "return-002", IdempotencyKey: "reject-001", ReviewedBy: "reviewer-002", ReviewNote: "damaged evidence missing"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.saved.Status != entities.ReturnRequestStatusRejected {
		t.Fatalf("expected status %s, got %s", entities.ReturnRequestStatusRejected, returns.saved.Status)
	}

	if returns.saved.ReviewedBy != "reviewer-002" {
		t.Fatalf("expected reviewer reviewer-002, got %s", returns.saved.ReviewedBy)
	}

	if returns.saved.ReviewNote != "damaged evidence missing" {
		t.Fatalf("expected review note to be saved, got %s", returns.saved.ReviewNote)
	}

	if idempotency.records["reject-return:reject-001"] != "return-002" {
		t.Fatal("expected idempotency record to be saved for rejected return")
	}
}

func TestAcceptReturnInteractorReusesSavedResultOnDuplicateKey(t *testing.T) {
	returns := &stubReturnRequestEditor{
		request: entities.ReturnRequest{
			ID:          "return-004",
			OrderID:     "order-001",
			Status:      entities.ReturnRequestStatusRefunded,
			RequestedBy: "customer-001",
			ReviewedBy:  "reviewer-001",
			ProcessedBy: "finance-001",
		},
	}
	orders := &stubOrderEditor{}
	restock := &stubInventoryRestock{}
	output := &stubAcceptReturnOutput{}
	refunds := &stubRefundGateway{}
	idempotency := &stubIdempotencyStore{
		records: map[string]string{
			"accept-return:accept-duplicate": "return-004",
		},
	}

	interactor := NewAcceptReturnInteractor(idempotency, orders, returns, stubReturnEligibilityPolicy{allowed: true}, refunds, restock, output)

	err := interactor.Execute(AcceptReturnInput{
		ReturnRequestID: "return-004",
		IdempotencyKey:  "accept-duplicate",
		ReviewedBy:      "reviewer-002",
		ProcessedBy:     "finance-002",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if refunds.calls != 0 {
		t.Fatalf("expected 0 refund calls on duplicate retry, got %d", refunds.calls)
	}

	if len(restock.items) != 0 {
		t.Fatalf("expected no restock items on duplicate retry, got %d", len(restock.items))
	}

	if returns.saved.ID != "" {
		t.Fatal("expected duplicate retry to avoid saving the return again")
	}

	if output.output.Status != entities.ReturnRequestStatusRefunded {
		t.Fatalf("expected presented status %s, got %s", entities.ReturnRequestStatusRefunded, output.output.Status)
	}
}

func TestRejectReturnInteractorReusesSavedResultOnDuplicateKey(t *testing.T) {
	returns := &stubReturnRequestEditor{
		request: entities.ReturnRequest{
			ID:         "return-005",
			OrderID:    "order-001",
			Status:     entities.ReturnRequestStatusRejected,
			ReviewedBy: "reviewer-002",
			ReviewNote: "damaged evidence missing",
		},
	}
	output := &stubRejectReturnOutput{}
	idempotency := &stubIdempotencyStore{
		records: map[string]string{
			"reject-return:reject-duplicate": "return-005",
		},
	}

	interactor := NewRejectReturnInteractor(idempotency, returns, output)

	err := interactor.Execute(RejectReturnInput{
		ReturnRequestID: "return-005",
		IdempotencyKey:  "reject-duplicate",
		ReviewedBy:      "reviewer-003",
		ReviewNote:      "new note",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.saved.ID != "" {
		t.Fatal("expected duplicate reject retry to avoid saving the return again")
	}

	if output.output.Status != entities.ReturnRequestStatusRejected {
		t.Fatalf("expected presented status %s, got %s", entities.ReturnRequestStatusRejected, output.output.Status)
	}
}

func TestAcceptReturnInteractorRequiresIdempotencyKey(t *testing.T) {
	interactor := NewAcceptReturnInteractor(&stubIdempotencyStore{}, &stubOrderEditor{}, &stubReturnRequestEditor{}, stubReturnEligibilityPolicy{allowed: true}, &stubRefundGateway{}, &stubInventoryRestock{}, &stubAcceptReturnOutput{})

	err := interactor.Execute(AcceptReturnInput{ReturnRequestID: "return-001", ReviewedBy: "reviewer-001", ProcessedBy: "finance-001"})
	if err != ErrIdempotencyKeyRequired {
		t.Fatalf("expected %v, got %v", ErrIdempotencyKeyRequired, err)
	}
}

func TestRejectReturnInteractorRequiresIdempotencyKey(t *testing.T) {
	interactor := NewRejectReturnInteractor(&stubIdempotencyStore{}, &stubReturnRequestEditor{}, &stubRejectReturnOutput{})

	err := interactor.Execute(RejectReturnInput{ReturnRequestID: "return-001", ReviewedBy: "reviewer-001", ReviewNote: "missing proof"})
	if err != ErrIdempotencyKeyRequired {
		t.Fatalf("expected %v, got %v", ErrIdempotencyKeyRequired, err)
	}
}
