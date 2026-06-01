package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubReturnRequestWriter struct {
	saved entities.ReturnRequest
}

func (g *stubReturnRequestWriter) Save(request entities.ReturnRequest) error {
	g.saved = request
	return nil
}

type stubRefundGateway struct {
	err error
}

func (g stubRefundGateway) Refund(order entities.Order) error {
	return g.err
}

type stubInventoryRestock struct {
	items []entities.InventoryReservationItem
	err   error
}

func (r *stubInventoryRestock) Restock(items []entities.InventoryReservationItem) error {
	if r.err != nil {
		return r.err
	}

	r.items = items
	return nil
}

type stubRequestReturnOutput struct {
	output RequestReturnOutput
}

func (o *stubRequestReturnOutput) Present(output RequestReturnOutput) error {
	o.output = output
	return nil
}

func TestRequestReturnInteractorCreatesRefundedReturnForShippedOrder(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusShipped,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", ProductName: "Office Chair", Quantity: 2},
			},
		},
	}
	returns := &stubReturnRequestWriter{}
	restock := &stubInventoryRestock{}
	output := &stubRequestReturnOutput{}

	interactor := NewRequestReturnInteractor(orders, returns, stubRefundGateway{}, restock, output)

	err := interactor.Execute(RequestReturnInput{OrderID: "order-001", Reason: "damaged item"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if returns.saved.OrderID != "order-001" {
		t.Fatalf("expected order id order-001, got %s", returns.saved.OrderID)
	}

	if returns.saved.Status != entities.ReturnRequestStatusRefunded {
		t.Fatalf("expected status %s, got %s", entities.ReturnRequestStatusRefunded, returns.saved.Status)
	}

	if len(restock.items) != 1 {
		t.Fatalf("expected 1 restock item, got %d", len(restock.items))
	}

	if restock.items[0].SKU != "CHAIR-001" {
		t.Fatalf("expected restocked sku CHAIR-001, got %s", restock.items[0].SKU)
	}
}

func TestRequestReturnInteractorRejectsNonShippedOrder(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-002",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPaid,
		},
	}
	returns := &stubReturnRequestWriter{}
	restock := &stubInventoryRestock{}
	output := &stubRequestReturnOutput{}

	interactor := NewRequestReturnInteractor(orders, returns, stubRefundGateway{}, restock, output)

	err := interactor.Execute(RequestReturnInput{OrderID: "order-002", Reason: "changed mind"})
	if err != entities.ErrOrderNotReturnable {
		t.Fatalf("expected %v, got %v", entities.ErrOrderNotReturnable, err)
	}
}
