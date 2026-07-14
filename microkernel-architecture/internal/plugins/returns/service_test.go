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

func TestRequestReturnStoresRequestedReturnWithoutRefundOrRestock(t *testing.T) {
	repository := &stubRepository{}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	service := NewService(repository, stubReturnableOrderProvider{
		order: kernel.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			ShippedAt:  time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC),
			Lines: []kernel.ReturnableOrderLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000, ReturnWindowDays: 30},
			},
		},
	}, stubClock{now: time.Date(2026, 7, 10, 12, 0, 0, 0, time.UTC)}, stubEligibilityPolicy{allowed: true}, refunds, restock)

	result, err := service.RequestReturn(kernel.RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged item",
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
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, refunds, restock)

	result, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
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
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, refunds, restock)

	_, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
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
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: true}, refunds, restock)

	result, err := service.RejectReturn(kernel.RejectReturnCommand{
		ReturnRequestID: "return-001",
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

	if refunds.orderID != "" || refunds.amount != 0 {
		t.Fatalf("expected no refund during rejection, got %s %d", refunds.orderID, refunds.amount)
	}

	if len(restock.items) != 0 {
		t.Fatalf("expected no restock during rejection, got %d items", len(restock.items))
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
	service := NewService(repository, stubReturnableOrderProvider{}, stubClock{}, stubEligibilityPolicy{allowed: false}, refunds, restock)

	result, err := service.AcceptReturn(kernel.AcceptReturnCommand{
		ReturnRequestID: "return-001",
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

	if refunds.orderID != "" || refunds.amount != 0 {
		t.Fatalf("expected no refund when policy blocks, got %s %d", refunds.orderID, refunds.amount)
	}

	if len(restock.items) != 0 {
		t.Fatalf("expected no restock when policy blocks, got %d items", len(restock.items))
	}
}
