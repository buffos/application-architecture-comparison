package returns

import (
	"testing"

	"modular-monolith/internal/modules/orders"
	"modular-monolith/internal/modules/payments"
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

func (s *stubRefunder) Refund(request payments.RefundRequest) error {
	if s.err != nil {
		return s.err
	}

	s.request = request
	return nil
}

func TestRequestReturnRefundsAndStoresReturnRequest(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	service := NewService(repository, stubOrderSource{
		order: orders.ReturnableOrder{
			OrderID:    "order-001",
			CustomerID: "customer-001",
			Lines: []orders.ReturnableOrderLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
			},
		},
	}, refunder)

	result, err := service.RequestReturn(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged item",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != ReturnRequestStatusRefunded {
		t.Fatalf("expected %s, got %s", ReturnRequestStatusRefunded, result.Status)
	}

	if refunder.request.Amount != 30000 {
		t.Fatalf("expected refund amount 30000, got %d", refunder.request.Amount)
	}
}

func TestRequestReturnRejectsNonReturnableOrder(t *testing.T) {
	repository := &stubRepository{}
	refunder := &stubRefunder{}
	service := NewService(repository, stubOrderSource{
		err: orders.ErrOrderNotReturnable,
	}, refunder)

	_, err := service.RequestReturn(RequestReturnCommand{
		OrderID: "order-001",
		Reason:  "damaged item",
	})
	if err != orders.ErrOrderNotReturnable {
		t.Fatalf("expected %v, got %v", orders.ErrOrderNotReturnable, err)
	}
}
