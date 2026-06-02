package orders

import (
	"testing"

	"modular-monolith/internal/modules/quotes"
)

type stubOrderRepository struct {
	saved Order
}

func (r *stubOrderRepository) Save(order Order) error {
	r.saved = order
	return nil
}

func (r *stubOrderRepository) FindByID(id string) (Order, error) {
	return r.saved, nil
}

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

func TestConvertQuoteToOrderCreatesPendingPaymentOrder(t *testing.T) {
	orders := &stubOrderRepository{}
	service := NewService(orders, stubApprovedQuoteSource{
		quote: quotes.ApprovedQuote{
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Lines: []quotes.ApprovedQuoteLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
			},
		},
	})

	result, err := service.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != OrderStatusPendingPayment {
		t.Fatalf("expected status %s, got %s", OrderStatusPendingPayment, result.Status)
	}

	if orders.saved.QuoteID != "quote-001" {
		t.Fatalf("expected quote-001, got %s", orders.saved.QuoteID)
	}
}

func TestConvertQuoteToOrderRejectsNonApprovedQuote(t *testing.T) {
	orders := &stubOrderRepository{}
	service := NewService(orders, stubApprovedQuoteSource{
		err: quotes.ErrQuoteNotConvertible,
	})

	_, err := service.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != quotes.ErrQuoteNotConvertible {
		t.Fatalf("expected %v, got %v", quotes.ErrQuoteNotConvertible, err)
	}
}
