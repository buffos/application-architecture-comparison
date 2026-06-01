package application

import (
	"testing"

	"onion-architecture/internal/domain"
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
			Lines: []domain.OrderLine{
				{
					ProductSKU: "sku-002",
					Quantity:   2,
				},
			},
		},
	}
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:      "return-001",
			OrderID: "order-001",
			Status:  domain.ReturnRequestStatusRequested,
			Reason:  "damaged on arrival",
		},
	}
	restock := &stubInventoryRestock{}

	service := NewAcceptReturnService(orders, returns, stubReturnEligibilityPolicy{eligible: true}, stubRefundGateway{}, restock)

	result, err := service.Execute(AcceptReturnCommand{ReturnRequestID: "return-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRefunded {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRefunded, result.Status)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected one restock item, got %d", len(restock.items))
	}
}

func TestRejectReturnServiceRejectsRequestedReturn(t *testing.T) {
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:      "return-001",
			OrderID: "order-001",
			Status:  domain.ReturnRequestStatusRequested,
		},
	}

	service := NewRejectReturnService(returns)

	result, err := service.Execute(RejectReturnCommand{ReturnRequestID: "return-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.ReturnRequestStatusRejected {
		t.Fatalf("expected status %s, got %s", domain.ReturnRequestStatusRejected, result.Status)
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
			Lines: []domain.OrderLine{
				{
					ProductSKU: "sku-002",
					Quantity:   2,
				},
			},
		},
	}
	returns := &stubReturnRequestStore{
		found: domain.ReturnRequest{
			ID:      "return-001",
			OrderID: "order-001",
			Status:  domain.ReturnRequestStatusRequested,
			Reason:  "outside return window",
		},
	}
	restock := &stubInventoryRestock{}

	service := NewAcceptReturnService(orders, returns, stubReturnEligibilityPolicy{eligible: false}, stubRefundGateway{}, restock)

	result, err := service.Execute(AcceptReturnCommand{ReturnRequestID: "return-001"})
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
