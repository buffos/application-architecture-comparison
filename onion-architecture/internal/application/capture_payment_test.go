package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubOrderRepository struct {
	order domain.Order
	err   error
	saved domain.Order
}

func (r *stubOrderRepository) FindByID(id string) (domain.Order, error) {
	if r.err != nil {
		return domain.Order{}, r.err
	}

	return r.order, nil
}

func (r *stubOrderRepository) Save(order domain.Order) error {
	r.saved = order
	return nil
}

type stubPaymentGateway struct {
	outcome string
	err     error
}

func (g stubPaymentGateway) Capture(order domain.Order) (string, error) {
	return g.outcome, g.err
}

func TestCapturePaymentServiceMarksPendingOrderPaid(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPendingPayment,
			Lines: []domain.OrderLine{
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

	service := NewCapturePaymentService(orders, stubPaymentGateway{outcome: PaymentCaptureOutcomeApproved})

	result, err := service.Execute(CapturePaymentCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.OrderStatusPaid {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusPaid, result.Status)
	}

	if orders.saved.Status != domain.OrderStatusPaid {
		t.Fatalf("expected saved status %s, got %s", domain.OrderStatusPaid, orders.saved.Status)
	}
}

func TestCapturePaymentServiceRejectsAlreadyPaidOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPaid,
		},
	}

	service := NewCapturePaymentService(orders, stubPaymentGateway{outcome: PaymentCaptureOutcomeApproved})

	_, err := service.Execute(CapturePaymentCommand{OrderID: "order-001"})
	if err != domain.ErrOrderNotPayable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotPayable, err)
	}
}

func TestCapturePaymentServiceMovesOrderToPaymentReview(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPendingPayment,
		},
	}

	service := NewCapturePaymentService(orders, stubPaymentGateway{outcome: PaymentCaptureOutcomeReview})

	result, err := service.Execute(CapturePaymentCommand{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Status != domain.OrderStatusPaymentReview {
		t.Fatalf("expected status %s, got %s", domain.OrderStatusPaymentReview, result.Status)
	}

	if orders.saved.Status != domain.OrderStatusPaymentReview {
		t.Fatalf("expected saved status %s, got %s", domain.OrderStatusPaymentReview, orders.saved.Status)
	}
}
