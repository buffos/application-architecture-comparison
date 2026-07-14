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

type stubInventoryRelease struct {
	err error
}

func (r stubInventoryRelease) Release(items []kernel.InventoryReservationItem) error {
	return r.err
}

type stubPaymentCapture struct {
	err error
}

func (p stubPaymentCapture) Capture(orderID string, amount int) error {
	return p.err
}

type stubShipmentCreation struct {
	result kernel.ShipmentCreationResult
	err    error
}

func (s stubShipmentCreation) CreateShipment(record kernel.CreateShipmentRecord) (kernel.ShipmentCreationResult, error) {
	return s.result, s.err
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
	}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{})

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
	}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{})

	_, err := service.ConvertQuoteToOrder(kernel.ConvertQuoteToOrderCommand{
		QuoteID: "quote-001",
	})
	if err != kernel.ErrPluginAlreadyRegistered {
		t.Fatalf("expected reservation error to propagate, got %v", err)
	}
}

func TestCapturePayment(t *testing.T) {
	repository := &stubRepository{
		saved: Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPendingPayment,
			Lines: []OrderLine{
				{
					ProductSKU:      "sku-001",
					ProductName:     "Desk",
					ProductCategory: "Standard",
					Quantity:        2,
					UnitPrice:       15000,
				},
			},
		},
	}
	service := NewService(repository, stubApprovedQuoteProvider{}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{})

	result, err := service.CapturePayment(kernel.CapturePaymentCommand{
		OrderID: "order-001",
	})
	if err != nil {
		t.Fatalf("expected capture payment to succeed, got %v", err)
	}

	if result.Status != OrderStatusPaid {
		t.Fatalf("expected paid status, got %s", result.Status)
	}
}

func TestCapturePaymentRejectsNonPayableOrder(t *testing.T) {
	repository := &stubRepository{
		saved: Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPaid,
		},
	}
	service := NewService(repository, stubApprovedQuoteProvider{}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{})

	_, err := service.CapturePayment(kernel.CapturePaymentCommand{
		OrderID: "order-001",
	})
	if err != ErrOrderNotPayable {
		t.Fatalf("expected not payable error, got %v", err)
	}
}

func TestCreateShipment(t *testing.T) {
	repository := &stubRepository{
		saved: Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPaid,
			Lines: []OrderLine{
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
	service := NewService(repository, stubApprovedQuoteProvider{}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{
		result: kernel.ShipmentCreationResult{
			ShipmentID: "shipment-001",
			OrderID:    "order-001",
			CustomerID: "customer-001",
			LineCount:  1,
		},
	})

	result, err := service.CreateShipment(kernel.CreateShipmentCommand{
		OrderID: "order-001",
	})
	if err != nil {
		t.Fatalf("expected create shipment to succeed, got %v", err)
	}

	if result.Status != OrderStatusShipped {
		t.Fatalf("expected shipped status, got %s", result.Status)
	}
}

func TestCreateShipmentRejectsNonShippableOrder(t *testing.T) {
	repository := &stubRepository{
		saved: Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPendingPayment,
		},
	}
	service := NewService(repository, stubApprovedQuoteProvider{}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{
		result: kernel.ShipmentCreationResult{
			ShipmentID: "shipment-001",
			OrderID:    "order-001",
			CustomerID: "customer-001",
			LineCount:  0,
		},
	})

	_, err := service.CreateShipment(kernel.CreateShipmentCommand{
		OrderID: "order-001",
	})
	if err != ErrOrderNotShippable {
		t.Fatalf("expected not shippable error, got %v", err)
	}
}

func TestCancelOrder(t *testing.T) {
	repository := &stubRepository{
		saved: Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPendingPayment,
			Lines: []OrderLine{
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
	service := NewService(repository, stubApprovedQuoteProvider{}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{})

	result, err := service.CancelOrder(kernel.CancelOrderCommand{
		OrderID: "order-001",
	})
	if err != nil {
		t.Fatalf("expected cancel order to succeed, got %v", err)
	}

	if result.Status != OrderStatusCancelled {
		t.Fatalf("expected cancelled status, got %s", result.Status)
	}
}

func TestCancelOrderRejectsShippedOrder(t *testing.T) {
	repository := &stubRepository{
		saved: Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     OrderStatusShipped,
		},
	}
	service := NewService(repository, stubApprovedQuoteProvider{}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{})

	_, err := service.CancelOrder(kernel.CancelOrderCommand{
		OrderID: "order-001",
	})
	if err != ErrOrderNotCancellable {
		t.Fatalf("expected not cancellable error, got %v", err)
	}
}

func TestGetReturnableOrder(t *testing.T) {
	repository := &stubRepository{
		saved: Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     OrderStatusShipped,
			Lines: []OrderLine{
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
	service := NewService(repository, stubApprovedQuoteProvider{}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{})

	result, err := service.GetReturnableOrder("order-001")
	if err != nil {
		t.Fatalf("expected returnable order lookup to succeed, got %v", err)
	}

	if result.OrderID != "order-001" {
		t.Fatalf("expected order id order-001, got %s", result.OrderID)
	}
}

func TestGetReturnableOrderRejectsNonShippedOrder(t *testing.T) {
	repository := &stubRepository{
		saved: Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     OrderStatusPaid,
		},
	}
	service := NewService(repository, stubApprovedQuoteProvider{}, stubInventoryReservation{}, stubInventoryRelease{}, stubPaymentCapture{}, stubShipmentCreation{})

	_, err := service.GetReturnableOrder("order-001")
	if err != ErrOrderNotReturnable {
		t.Fatalf("expected not returnable error, got %v", err)
	}
}
