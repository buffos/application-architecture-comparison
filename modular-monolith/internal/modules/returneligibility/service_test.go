package returneligibility

import "testing"

func TestAllowsEligibleReturn(t *testing.T) {
	service := NewService()

	if !service.Allows(ReviewRequest{Reason: "damaged item"}) {
		t.Fatalf("expected damaged item return to be allowed")
	}
}

func TestRejectsOutsideReturnWindow(t *testing.T) {
	service := NewService()

	if service.Allows(ReviewRequest{Reason: "outside return window"}) {
		t.Fatalf("expected outside return window to be rejected")
	}
}
