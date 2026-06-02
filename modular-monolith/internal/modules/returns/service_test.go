package returns

import (
	"testing"

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

func TestRequestReturnStoresRequestedReturn(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	service := NewService(repository, stubOrderSource{
		order: orders.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			Lines: []orders.ReturnableOrderLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
			},
		},
	}, stubEligibility{allows: true}, restocker, refunder)

	result, err := service.RequestReturn(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged item",
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
}

func TestRequestReturnRejectsNonReturnableOrder(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	service := NewService(repository, stubOrderSource{
		err: orders.ErrOrderNotReturnable,
	}, stubEligibility{allows: true}, restocker, refunder)

	_, err := service.RequestReturn(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged item",
	})
	if err != orders.ErrOrderNotReturnable {
		t.Fatalf("expected %v, got %v", orders.ErrOrderNotReturnable, err)
	}
}

func TestRequestReturnStopsWhenRestockFails(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{err: inventory.ErrStockNotFound}
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: true}, restocker, refunder)
	repository.saved = NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
		},
	}, "damaged item")

	_, err := service.AcceptReturn(ReviewReturnCommand{
		ReturnRequestID: repository.saved.ID,
	})
	if err != inventory.ErrStockNotFound {
		t.Fatalf("expected %v, got %v", inventory.ErrStockNotFound, err)
	}
}

func TestAcceptReturnRefundsRestocksAndStoresUpdatedStatus(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	repository.saved = NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
		},
	}, "damaged item")
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: true}, restocker, refunder)

	result, err := service.AcceptReturn(ReviewReturnCommand{ReturnRequestID: repository.saved.ID})
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
}

func TestRejectReturnStoresRejectedStatus(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	repository.saved = NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
		},
	}, "damaged item")
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: true}, restocker, refunder)

	result, err := service.RejectReturn(ReviewReturnCommand{ReturnRequestID: repository.saved.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != ReturnRequestStatusRejected {
		t.Fatalf("expected %s, got %s", ReturnRequestStatusRejected, result.Status)
	}

	if refunder.request.Amount != 0 {
		t.Fatalf("expected no refund on rejection, got %d", refunder.request.Amount)
	}
}

func TestAcceptReturnRejectsWhenPolicyBlocksEligibility(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	restocker := &stubRestocker{}
	repository.saved = NewRequestedReturnRequest(ReturnableOrder{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Lines: []ReturnableOrderLine{
			{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
		},
	}, "outside return window")
	service := NewService(repository, stubOrderSource{}, stubEligibility{allows: false}, restocker, refunder)

	result, err := service.AcceptReturn(ReviewReturnCommand{ReturnRequestID: repository.saved.ID})
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
}
