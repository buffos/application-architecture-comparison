package returns

import "testing"

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
