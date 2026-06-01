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
	output := &stubConvertQuoteToOrderOutput{}

	interactor := NewConvertQuoteToOrderInteractor(quotes, orders, output)

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
	output := &stubConvertQuoteToOrderOutput{}

	interactor := NewConvertQuoteToOrderInteractor(quotes, orders, output)

	err := interactor.Execute(ConvertQuoteToOrderInput{QuoteID: "quote-002"})
	if err != entities.ErrQuoteNotConvertible {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteNotConvertible, err)
	}
}
