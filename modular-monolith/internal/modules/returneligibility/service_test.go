package returneligibility

import (
	"testing"
	"time"
)

func TestAllowsEligibleReturn(t *testing.T) {
	service := NewService()

	shippedAt := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	requestedAt := shippedAt.AddDate(0, 0, 10)
	if !service.Allows(ReviewRequest{
		Reason:      "damaged item",
		ShippedAt:   shippedAt,
		RequestedAt: requestedAt,
		Lines:       []ReviewLine{{ReturnWindowDays: 30}},
	}) {
		t.Fatalf("expected damaged item return to be allowed")
	}
}

func TestRejectsOutsideReturnWindow(t *testing.T) {
	service := NewService()

	shippedAt := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	requestedAt := shippedAt.AddDate(0, 0, 31)
	if service.Allows(ReviewRequest{
		Reason:      "damaged item",
		ShippedAt:   shippedAt,
		RequestedAt: requestedAt,
		Lines:       []ReviewLine{{ReturnWindowDays: 30}},
	}) {
		t.Fatalf("expected out-of-window return to be rejected")
	}
}

func TestRejectsExplicitOutsideReturnWindowReason(t *testing.T) {
	service := NewService()

	shippedAt := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
	requestedAt := shippedAt.AddDate(0, 0, 10)
	if service.Allows(ReviewRequest{
		Reason:      "outside return window",
		ShippedAt:   shippedAt,
		RequestedAt: requestedAt,
		Lines:       []ReviewLine{{ReturnWindowDays: 30}},
	}) {
		t.Fatalf("expected outside return window to be rejected")
	}
}
