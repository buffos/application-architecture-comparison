package returns

import (
	"errors"
	"testing"

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

func TestRequestReturn(t *testing.T) {
	repository := &stubRepository{}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{}
	service := NewService(repository, stubReturnableOrderProvider{
		order: kernel.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			Lines: []kernel.ReturnableOrderLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000},
			},
		},
	}, refunds, restock)

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

	if refunds.orderID != "order-001" || refunds.amount != 45000 {
		t.Fatalf("expected refund for order-001 amount 45000, got %s %d", refunds.orderID, refunds.amount)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected 1 restock item, got %d", len(restock.items))
	}

	if restock.items[0].ProductSKU != "sku-002" || restock.items[0].Quantity != 1 {
		t.Fatalf("unexpected restock item %+v", restock.items[0])
	}
}

func TestRequestReturnStopsWhenRestockFails(t *testing.T) {
	repository := &stubRepository{}
	refunds := &stubPaymentRefund{}
	restock := &stubInventoryRestock{err: errors.New("restock failed")}
	service := NewService(repository, stubReturnableOrderProvider{
		order: kernel.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			Lines: []kernel.ReturnableOrderLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000},
			},
		},
	}, refunds, restock)

	_, err := service.RequestReturn(kernel.RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged item",
	})
	if err == nil || err.Error() != "restock failed" {
		t.Fatalf("expected restock failure, got %v", err)
	}

	if repository.saved.ID != "" {
		t.Fatalf("expected return request not to be saved when restock fails")
	}

	if refunds.orderID != "order-001" {
		t.Fatalf("expected refund to be attempted before restock failure, got %s", refunds.orderID)
	}
}
