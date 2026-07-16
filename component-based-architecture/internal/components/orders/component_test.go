package orders

import (
	"errors"
	"testing"

	"component-based-architecture/internal/components/quotes"
)

type stubApprovedQuoteSource struct {
	quote quotes.ApprovedQuote
	err   error
}

func (s stubApprovedQuoteSource) GetApprovedQuoteForOrder(quoteID string) (quotes.ApprovedQuote, error) {
	if s.err != nil {
		return quotes.ApprovedQuote{}, s.err
	}
	return s.quote, nil
}

func TestConvertQuoteToOrderCreatesOrderSnapshot(t *testing.T) {
	component := NewComponent(stubApprovedQuoteSource{quote: quotes.ApprovedQuote{
		QuoteID: "quote-001", CustomerID: "customer-001",
		Lines: []quotes.ApprovedQuoteLine{{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000}},
	}})

	result, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("convert quote: %v", err)
	}
	if result.OrderID != "order-001" {
		t.Fatalf("expected order-001, got %s", result.OrderID)
	}
	if result.Status != OrderStatusPendingPayment {
		t.Fatalf("expected %s, got %s", OrderStatusPendingPayment, result.Status)
	}
	if result.LineCount != 1 {
		t.Fatalf("expected one line, got %d", result.LineCount)
	}
}

func TestConvertQuoteToOrderPropagatesNonApprovedQuoteError(t *testing.T) {
	component := NewComponent(stubApprovedQuoteSource{err: quotes.ErrQuoteNotConvertible})

	_, err := component.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if !errors.Is(err, quotes.ErrQuoteNotConvertible) {
		t.Fatalf("expected %v, got %v", quotes.ErrQuoteNotConvertible, err)
	}
}
