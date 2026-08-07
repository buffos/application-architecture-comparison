package returns

import (
	"testing"
	"time"
)

func TestReturnEligibilityServiceEvaluatesRealReturnWindows(t *testing.T) {
	service := NewReturnEligibilityService()
	shippedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	inside := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	if decision := service.EvaluateWindow(shippedAt, inside, []EligibilityLine{{Category: ProductCategoryStandard, ReturnWindowDays: 30}}); !decision.Eligible {
		t.Fatalf("inside-window decision = %+v", decision)
	}

	expired := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	if decision := service.EvaluateWindow(shippedAt, expired, []EligibilityLine{{Category: ProductCategoryStandard, ReturnWindowDays: 30}}); decision.Eligible || decision.Reason != "return window has expired" {
		t.Fatalf("expired decision = %+v", decision)
	}
}

func TestClearancePrecedesWindowEvaluation(t *testing.T) {
	service := NewReturnEligibilityService()
	shippedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	requestedAt := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	decision := service.EvaluateWindow(shippedAt, requestedAt, []EligibilityLine{{Category: ProductCategoryClearance, ReturnWindowDays: 30}})
	if decision.Eligible || decision.Reason != "clearance products are not returnable" {
		t.Fatalf("clearance decision = %+v", decision)
	}
}
