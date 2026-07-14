package returneligibility

import (
	"testing"
	"time"

	"microkernel-architecture/internal/kernel"
)

func TestAllowsInWindowReturn(t *testing.T) {
	service := NewService()

	shippedAt := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	requestedAt := shippedAt.AddDate(0, 0, 10)
	if !service.Allows(kernel.ReturnEligibilityReview{
		Reason:      "damaged item",
		ShippedAt:   shippedAt,
		RequestedAt: requestedAt,
		Lines:       []kernel.ReturnEligibilityLine{{ReturnWindowDays: 30}},
	}) {
		t.Fatalf("expected in-window return to be allowed")
	}
}

func TestRejectsOutOfWindowReturn(t *testing.T) {
	service := NewService()

	shippedAt := time.Date(2026, 7, 1, 12, 0, 0, 0, time.UTC)
	requestedAt := shippedAt.AddDate(0, 0, 31)
	if service.Allows(kernel.ReturnEligibilityReview{
		Reason:      "damaged item",
		ShippedAt:   shippedAt,
		RequestedAt: requestedAt,
		Lines:       []kernel.ReturnEligibilityLine{{ReturnWindowDays: 30}},
	}) {
		t.Fatalf("expected out-of-window return to be rejected")
	}
}
