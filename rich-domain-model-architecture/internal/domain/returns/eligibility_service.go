package returns

import "time"

type EligibilityLine struct {
	Category         ProductCategory
	ReturnWindowDays int
}

type EligibilityDecision struct {
	Eligible bool
	Reason   string
}

type ReturnEligibilityService struct{}

func NewReturnEligibilityService() ReturnEligibilityService {
	return ReturnEligibilityService{}
}

func (ReturnEligibilityService) Evaluate(lines []EligibilityLine) EligibilityDecision {
	for _, line := range lines {
		if line.Category == ProductCategoryClearance {
			return EligibilityDecision{Eligible: false, Reason: "clearance products are not returnable"}
		}
	}
	return EligibilityDecision{Eligible: true}
}

func (service ReturnEligibilityService) EvaluateWindow(shippedAt, requestedAt time.Time, lines []EligibilityLine) EligibilityDecision {
	decision := service.Evaluate(lines)
	if !decision.Eligible {
		return decision
	}
	for _, line := range lines {
		if line.ReturnWindowDays <= 0 || requestedAt.After(shippedAt.AddDate(0, 0, line.ReturnWindowDays)) {
			return EligibilityDecision{Eligible: false, Reason: "return window has expired"}
		}
	}
	return decision
}
