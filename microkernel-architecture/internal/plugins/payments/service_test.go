package payments

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

func TestServiceCapturesGatewayOutcome(t *testing.T) {
	service := NewService(NewManualReviewGateway())

	result, err := service.Capture("order-001", 45000)
	if err != nil {
		t.Fatalf("expected capture to succeed, got %v", err)
	}

	if result.Outcome != kernel.PaymentCaptureOutcomeReview {
		t.Fatalf("expected review outcome, got %s", result.Outcome)
	}
}
