package returns

import (
	"testing"
	"time"
)

func TestReturnEligibilityRejectsClearance(t *testing.T) {
	decision := NewReturnEligibilityService().Evaluate([]EligibilityLine{{Category: ProductCategoryStandard}, {Category: ProductCategoryClearance}})
	if decision.Eligible || decision.Reason == "" {
		t.Fatalf("unexpected decision %+v", decision)
	}
}

func TestReturnEligibilityAllowsReturnableCategories(t *testing.T) {
	decision := NewReturnEligibilityService().Evaluate([]EligibilityLine{{Category: ProductCategoryStandard}, {Category: ProductCategoryCustomBuild}})
	if !decision.Eligible || decision.Reason != "" {
		t.Fatalf("unexpected decision %+v", decision)
	}
}

func TestReturnEligibilityEnforcesReturnWindow(t *testing.T) {
	shippedAt := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	inside := NewReturnEligibilityService().EvaluateWindow(shippedAt, shippedAt.AddDate(0, 0, 30), []EligibilityLine{{Category: ProductCategoryStandard, ReturnWindowDays: 30}})
	if !inside.Eligible {
		t.Fatalf("inside-window decision %+v", inside)
	}
	outside := NewReturnEligibilityService().EvaluateWindow(shippedAt, shippedAt.AddDate(0, 0, 31), []EligibilityLine{{Category: ProductCategoryStandard, ReturnWindowDays: 30}})
	if outside.Eligible || outside.Reason == "" {
		t.Fatalf("outside-window decision %+v", outside)
	}
}

func TestReturnEligibilityKeepsClearancePrecedence(t *testing.T) {
	decision := NewReturnEligibilityService().EvaluateWindow(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC), []EligibilityLine{{Category: ProductCategoryClearance, ReturnWindowDays: 30}})
	if decision.Eligible || decision.Reason != "clearance products are not returnable" {
		t.Fatalf("unexpected decision %+v", decision)
	}
}
