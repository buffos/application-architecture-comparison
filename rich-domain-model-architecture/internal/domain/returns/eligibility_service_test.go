package returns

import "testing"

func TestReturnEligibilityServiceEvaluatesProductCategories(t *testing.T) {
	service := NewReturnEligibilityService()
	if decision := service.Evaluate([]EligibilityLine{{Category: ProductCategoryStandard}, {Category: ProductCategoryCustomBuild}}); !decision.Eligible || decision.Reason != "" {
		t.Fatalf("eligible decision = %+v", decision)
	}
	decision := service.Evaluate([]EligibilityLine{{Category: ProductCategoryStandard}, {Category: ProductCategoryClearance}})
	if decision.Eligible || decision.Reason != "clearance products are not returnable" {
		t.Fatalf("clearance decision = %+v", decision)
	}
}

func TestReturnEligibilityServiceDoesNotChangeInputs(t *testing.T) {
	lines := []EligibilityLine{{Category: ProductCategoryStandard}}
	_ = NewReturnEligibilityService().Evaluate(lines)
	if lines[0].Category != ProductCategoryStandard {
		t.Fatal("eligibility evaluation changed the input")
	}
}
