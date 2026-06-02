package orders

import (
	"testing"

	"modular-monolith/internal/modules/inventory"
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

type stubInventoryReserver struct {
	reserved []inventory.ReservationItem
	err      error
}

func (s *stubInventoryReserver) Reserve(items []inventory.ReservationItem) error {
	if s.err != nil {
		return s.err
	}

	s.reserved = append([]inventory.ReservationItem(nil), items...)
	return nil
}

func (s stubApprovedQuoteSource) GetApprovedQuoteForOrder(quoteID string) (quotes.ApprovedQuote, error) {
	if s.err != nil {
		return quotes.ApprovedQuote{}, s.err
	}

	return s.quote, nil
}

func TestConvertQuoteToOrderCreatesPendingPaymentOrder(t *testing.T) {
	orders := &stubOrderRepository{}
	reserver := &stubInventoryReserver{}
	service := NewService(orders, stubApprovedQuoteSource{
		quote: quotes.ApprovedQuote{
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Lines: []quotes.ApprovedQuoteLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
			},
		},
	}, reserver)

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

	if len(reserver.reserved) != 1 || reserver.reserved[0].Quantity != 2 {
		t.Fatalf("expected reservation for quantity 2, got %+v", reserver.reserved)
	}
}

func TestConvertQuoteToOrderRejectsNonApprovedQuote(t *testing.T) {
	orders := &stubOrderRepository{}
	reserver := &stubInventoryReserver{}
	service := NewService(orders, stubApprovedQuoteSource{
		err: quotes.ErrQuoteNotConvertible,
	}, reserver)

	_, err := service.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != quotes.ErrQuoteNotConvertible {
		t.Fatalf("expected %v, got %v", quotes.ErrQuoteNotConvertible, err)
	}
}

func TestConvertQuoteToOrderStopsWhenReservationFails(t *testing.T) {
	orders := &stubOrderRepository{}
	reserver := &stubInventoryReserver{err: inventory.ErrInsufficientStock}
	service := NewService(orders, stubApprovedQuoteSource{
		quote: quotes.ApprovedQuote{
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Lines: []quotes.ApprovedQuoteLine{
				{ProductSKU: "sku-001", ProductName: "Desk", ProductCategory: "Standard", Quantity: 2, UnitPrice: 15000},
			},
		},
	}, reserver)

	_, err := service.ConvertQuoteToOrder(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != inventory.ErrInsufficientStock {
		t.Fatalf("expected %v, got %v", inventory.ErrInsufficientStock, err)
	}

	if orders.saved.ID != "" {
		t.Fatalf("expected order not to be saved when reservation fails")
	}
}
