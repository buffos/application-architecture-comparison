package orders

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	saved Order
}

func (r *stubRepository) FindByID(id string) (Order, error) {
	if r.saved.ID == id {
		return r.saved, nil
	}

	return Order{}, ErrOrderNotFound
}

func (r *stubRepository) Save(order Order) error {
	r.saved = order
	return nil
}

type stubApprovedQuoteProvider struct {
	quote kernel.ApprovedQuote
	err   error
}

func (p stubApprovedQuoteProvider) GetApprovedQuoteForOrder(quoteID string) (kernel.ApprovedQuote, error) {
	return p.quote, p.err
}

type stubInventoryReservation struct {
	err error
}

func (r stubInventoryReservation) Reserve(items []kernel.InventoryReservationItem) error {
	return r.err
}

func TestConvertQuoteToOrder(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository, stubApprovedQuoteProvider{
		quote: kernel.ApprovedQuote{
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Lines: []kernel.ApprovedQuoteLine{
				{
					ProductSKU:      "sku-001",
					ProductName:     "Desk",
					ProductCategory: "Standard",
					Quantity:        2,
					UnitPrice:       15000,
				},
			},
		},
	}, stubInventoryReservation{})

	result, err := service.ConvertQuoteToOrder(kernel.ConvertQuoteToOrderCommand{
		QuoteID: "quote-001",
	})
	if err != nil {
		t.Fatalf("expected convert quote to order to succeed, got %v", err)
	}

	if result.QuoteID != "quote-001" {
		t.Fatalf("expected quote id quote-001, got %s", result.QuoteID)
	}

	if repository.saved.Status != OrderStatusPendingPayment {
		t.Fatalf("expected pending payment status, got %s", repository.saved.Status)
	}
}

func TestConvertQuoteToOrderRejectsReservationFailure(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository, stubApprovedQuoteProvider{
		quote: kernel.ApprovedQuote{
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Lines: []kernel.ApprovedQuoteLine{
				{
					ProductSKU:      "sku-001",
					ProductName:     "Desk",
					ProductCategory: "Standard",
					Quantity:        2,
					UnitPrice:       15000,
				},
			},
		},
	}, stubInventoryReservation{
		err: kernel.ErrPluginAlreadyRegistered,
	})

	_, err := service.ConvertQuoteToOrder(kernel.ConvertQuoteToOrderCommand{
		QuoteID: "quote-001",
	})
	if err != kernel.ErrPluginAlreadyRegistered {
		t.Fatalf("expected reservation error to propagate, got %v", err)
	}
}
