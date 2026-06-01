package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

func TestApprovePaymentReviewServiceMarksReviewedOrderPaid(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:         "order-001",
			QuoteID:    "quote-001",
			CustomerID: "customer-001",
			Status:     domain.OrderStatusPaymentReview,
		},
	}

	service := NewApprovePaymentReviewService(orders)

	result, err := service.Execute(ApprovePaymentReviewCommand{OrderID: "order-001"})
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

func TestApprovePaymentReviewServiceRejectsNonReviewedOrder(t *testing.T) {
	orders := &stubOrderRepository{
		order: domain.Order{
			ID:     "order-001",
			Status: domain.OrderStatusPendingPayment,
		},
	}

	service := NewApprovePaymentReviewService(orders)

	_, err := service.Execute(ApprovePaymentReviewCommand{OrderID: "order-001"})
	if err != domain.ErrOrderNotPaymentReviewable {
		t.Fatalf("expected %v, got %v", domain.ErrOrderNotPaymentReviewable, err)
	}
}
