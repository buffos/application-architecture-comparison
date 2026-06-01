package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubInventoryRelease struct {
	items []entities.InventoryReservationItem
	err   error
}

func (r *stubInventoryRelease) Release(items []entities.InventoryReservationItem) error {
	if r.err != nil {
		return r.err
	}

	r.items = items
	return nil
}

type stubCancelOrderOutput struct {
	output CancelOrderOutput
}

func (o *stubCancelOrderOutput) Present(output CancelOrderOutput) error {
	o.output = output
	return nil
}

func TestCancelOrderInteractorCancelsUnshippedOrderAndReleasesStock(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-001",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPaid,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}
	inventory := &stubInventoryRelease{}
	output := &stubCancelOrderOutput{}

	interactor := NewCancelOrderInteractor(orders, inventory, output)

	err := interactor.Execute(CancelOrderInput{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if orders.saved.Status != entities.OrderStatusCancelled {
		t.Fatalf("expected saved status %s, got %s", entities.OrderStatusCancelled, orders.saved.Status)
	}

	if len(inventory.items) != 1 {
		t.Fatalf("expected 1 released item, got %d", len(inventory.items))
	}
}

func TestCancelOrderInteractorRejectsShippedOrder(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-002",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusShipped,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}
	inventory := &stubInventoryRelease{}
	output := &stubCancelOrderOutput{}

	interactor := NewCancelOrderInteractor(orders, inventory, output)

	err := interactor.Execute(CancelOrderInput{OrderID: "order-002"})
	if err != entities.ErrQuoteCannotTransition {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotTransition, err)
	}
}

func TestCancelOrderInteractorRejectsPartiallyShippedOrder(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-003",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPartiallyShipped,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2, ShippedQuantity: 1},
			},
		},
	}
	inventory := &stubInventoryRelease{}
	output := &stubCancelOrderOutput{}

	interactor := NewCancelOrderInteractor(orders, inventory, output)

	err := interactor.Execute(CancelOrderInput{OrderID: "order-003"})
	if err != entities.ErrQuoteCannotTransition {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotTransition, err)
	}
}
