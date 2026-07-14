package returneligibility

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

func TestAllowsEligibleReason(t *testing.T) {
	service := NewService()

	if !service.Allows(kernel.ReturnEligibilityReview{Reason: "damaged item"}) {
		t.Fatalf("expected damaged item return to be allowed")
	}
}

func TestRejectsOutsideReturnWindowReason(t *testing.T) {
	service := NewService()

	if service.Allows(kernel.ReturnEligibilityReview{Reason: "outside return window"}) {
		t.Fatalf("expected outside return window to be rejected")
	}
}
