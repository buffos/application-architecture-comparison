package returns

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
