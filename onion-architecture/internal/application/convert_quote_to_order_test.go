package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubOrderStore struct {
	saved domain.Order
}

func (s *stubOrderStore) Save(order domain.Order) error {
	s.saved = order
	return nil
}

func TestConvertQuoteToOrderServiceCreatesOrderFromApprovedQuote(t *testing.T) {
	quotes := stubQuoteFinder{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusApproved,
			Lines: []domain.QuoteLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        1,
					UnitPrice:       45000,
				},
			},
		},
	}
	orders := &stubOrderStore{}

	service := NewConvertQuoteToOrderService(quotes, orders)

	result, err := service.Execute(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.OrderStatusPendingPayment {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusPendingPayment, result.Status)
	}

	if orders.saved.QuoteID != "quote-001" {
		t.Fatalf("expected saved quote id quote-001, got %s", orders.saved.QuoteID)
	}
}

func TestConvertQuoteToOrderServiceRejectsNonApprovedQuote(t *testing.T) {
	quotes := stubQuoteFinder{
		quote: domain.Quote{
			ID:         "quote-001",
			CustomerID: "customer-001",
			Status:     domain.QuoteStatusPendingApproval,
			Lines: []domain.QuoteLine{
				{
					ProductSKU:      "sku-002",
					ProductName:     "Custom Desk",
					ProductCategory: "CustomBuild",
					Quantity:        1,
					UnitPrice:       45000,
				},
			},
		},
	}
	orders := &stubOrderStore{}

	service := NewConvertQuoteToOrderService(quotes, orders)

	_, err := service.Execute(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != domain.ErrQuoteNotConvertible {
		t.Fatalf("expected %v, got %v", domain.ErrQuoteNotConvertible, err)
	}
}
