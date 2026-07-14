package returns

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	saved ReturnRequest
}

func (r *stubRepository) FindByID(id string) (ReturnRequest, error) {
	if r.saved.ID == id {
		return r.saved, nil
	}

	return ReturnRequest{}, ErrReturnRequestNotFound
}

func (r *stubRepository) Save(request ReturnRequest) error {
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
	err error
}

func (p stubPaymentRefund) Refund(orderID string, amount int) error {
	return p.err
}

func TestRequestReturn(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository, stubReturnableOrderProvider{
		order: kernel.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			Lines: []kernel.ReturnableOrderLine{
				{ProductSKU: "sku-002", Quantity: 1, UnitPrice: 45000},
			},
		},
	}, stubPaymentRefund{})

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
}
