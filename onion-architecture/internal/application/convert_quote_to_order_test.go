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

type stubInventoryReservation struct {
	items []domain.InventoryReservationItem
	err   error
}

func (s *stubInventoryReservation) Reserve(items []domain.InventoryReservationItem) error {
	if s.err != nil {
		return s.err
	}

	s.items = items
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
	inventory := &stubInventoryReservation{}

	service := NewConvertQuoteToOrderService(quotes, orders, inventory)

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

	if len(inventory.items) != 1 {
		t.Fatalf("expected one reservation item, got %d", len(inventory.items))
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
	inventory := &stubInventoryReservation{}

	service := NewConvertQuoteToOrderService(quotes, orders, inventory)

	_, err := service.Execute(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != domain.ErrQuoteNotConvertible {
		t.Fatalf("expected %v, got %v", domain.ErrQuoteNotConvertible, err)
	}
}

func TestConvertQuoteToOrderServiceRejectsWhenReservationFails(t *testing.T) {
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
					Quantity:        2,
					UnitPrice:       45000,
				},
			},
		},
	}
	orders := &stubOrderStore{}
	inventory := &stubInventoryReservation{err: domain.ErrInsufficientStock}

	service := NewConvertQuoteToOrderService(quotes, orders, inventory)

	_, err := service.Execute(ConvertQuoteToOrderCommand{QuoteID: "quote-001"})
	if err != domain.ErrInsufficientStock {
		t.Fatalf("expected %v, got %v", domain.ErrInsufficientStock, err)
	}

	if orders.saved.ID != "" {
		t.Fatalf("expected order not to be saved when reservation fails")
	}
}
