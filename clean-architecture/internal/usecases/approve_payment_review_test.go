package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubApprovePaymentReviewOutput struct {
	output ApprovePaymentReviewOutput
}

func (o *stubApprovePaymentReviewOutput) Present(output ApprovePaymentReviewOutput) error {
	o.output = output
	return nil
}

func TestApprovePaymentReviewInteractorMarksOrderPaid(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-010",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPaymentReview,
			Lines: []entities.OrderLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}
	output := &stubApprovePaymentReviewOutput{}

	interactor := NewApprovePaymentReviewInteractor(orders, output)

	err := interactor.Execute(ApprovePaymentReviewInput{OrderID: "order-010"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if orders.saved.Status != entities.OrderStatusPaid {
		t.Fatalf("expected saved status %s, got %s", entities.OrderStatusPaid, orders.saved.Status)
	}
}

func TestApprovePaymentReviewInteractorRejectsWrongOrderState(t *testing.T) {
	orders := &stubOrderEditor{
		order: entities.Order{
			ID:            "order-011",
			CustomerID:    "customer-001",
			SourceQuoteID: "quote-001",
			Status:        entities.OrderStatusPendingPayment,
		},
	}
	output := &stubApprovePaymentReviewOutput{}

	interactor := NewApprovePaymentReviewInteractor(orders, output)

	err := interactor.Execute(ApprovePaymentReviewInput{OrderID: "order-011"})
	if err != entities.ErrQuoteCannotTransition {
		t.Fatalf("expected %v, got %v", entities.ErrQuoteCannotTransition, err)
	}
}
