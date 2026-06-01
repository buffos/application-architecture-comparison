package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubOrderWriter struct {
	saved entities.Order
}

func (g *stubOrderWriter) Save(order entities.Order) error {
	g.saved = order
	return nil
}

type stubConvertQuoteToOrderOutput struct {
	output ConvertQuoteToOrderOutput
}

func (o *stubConvertQuoteToOrderOutput) Present(output ConvertQuoteToOrderOutput) error {
	o.output = output
	return nil
}

type stubInventoryReservation struct {
	items []entities.InventoryReservationItem
	err   error
}

func (r *stubInventoryReservation) Reserve(items []entities.InventoryReservationItem) error {
	if r.err != nil {
		return r.err
	}

	r.items = items
	return nil
}

func TestConvertQuoteToOrderInteractorCreatesOrderFromApprovedQuote(t *testing.T) {
	quotes := stubQuoteReader{
		quote: entities.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusApproved,
			Lines: []entities.QuoteLine{
				{SKU: "CHAIR-001", ProductName: "Office Chair", Quantity: 2, UnitPrice: 10000, LineTotal: 20000},
			},
		},
	}
	orders := &stubOrderWriter{}
	inventory := &stubInventoryReservation{}
	output := &stubConvertQuoteToOrderOutput{}

	interactor := NewConvertQuoteToOrderInteractor(quotes, orders, inventory, output)

	err := interactor.Execute(ConvertQuoteToOrderInput{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if orders.saved.SourceQuoteID != "quote-001" {
		t.Fatalf("expected source quote id quote-001, got %s", orders.saved.SourceQuoteID)
	}

	if orders.saved.Status != entities.OrderStatusPendingPayment {
		t.Fatalf("expected status %s, got %s", entities.OrderStatusPendingPayment, orders.saved.Status)
	}

	if len(inventory.items) != 1 {
		t.Fatalf("expected 1 reservation item, got %d", len(inventory.items))
	}

	if output.output.OrderID == "" {
		t.Fatal("expected presenter output to include order id")
	}
}

func TestConvertQuoteToOrderInteractorRejectsNonApprovedQuote(t *testing.T) {
	quotes := stubQuoteReader{
		quote: entities.Quote{
			ID:         "quote-002",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusPendingApproval,
			Lines: []entities.QuoteLine{
				{SKU: "DESK-001", Quantity: 1},
			},
		},
	}
	orders := &stubOrderWriter{}
	inventory := &stubInventoryReservation{}
	output := &stubConvertQuoteToOrderOutput{}

	interactor := NewConvertQuoteToOrderInteractor(quotes, orders, inventory, output)

	err := interactor.Execute(ConvertQuoteToOrderInput{QuoteID: "quote-002"})
	if err != entities.ErrQuoteNotConvertible {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteNotConvertible, err)
	}
}

func TestConvertQuoteToOrderInteractorFailsWhenReservationFails(t *testing.T) {
	quotes := stubQuoteReader{
		quote: entities.Quote{
			ID:         "quote-003",
			CustomerID: "customer-001",
			Status:     entities.QuoteStatusApproved,
			Lines: []entities.QuoteLine{
				{SKU: "CHAIR-001", ProductName: "Office Chair", Quantity: 5, UnitPrice: 10000, LineTotal: 50000},
			},
		},
	}
	orders := &stubOrderWriter{}
	inventory := &stubInventoryReservation{err: entities.ErrInsufficientInventory}
	output := &stubConvertQuoteToOrderOutput{}

	interactor := NewConvertQuoteToOrderInteractor(quotes, orders, inventory, output)

	err := interactor.Execute(ConvertQuoteToOrderInput{QuoteID: "quote-003"})
	if err != entities.ErrInsufficientInventory {
		t.Fatalf("expected %v, got %v", entities.ErrInsufficientInventory, err)
	}

	if orders.saved.ID != "" {
		t.Fatal("expected no order to be saved when reservation fails")
	}
}
