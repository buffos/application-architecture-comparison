package rules

import "rules-engine-architecture/internal/engine"

// PreferredDiscountEligibilityRule publishes eligibility for a future
// discount policy without choosing or applying a concrete discount amount.
type PreferredDiscountEligibilityRule struct{}

func (PreferredDiscountEligibilityRule) Name() string {
	return "Preferred Discount Eligibility Rule"
}

func (PreferredDiscountEligibilityRule) Priority() int {
	return 40
}

func (PreferredDiscountEligibilityRule) ConflictGroup() string {
	return ""
}

func (PreferredDiscountEligibilityRule) Evaluate(memory *engine.WorkingMemory) bool {
	return memory.Customer.Tier == "Preferred"
}

func (PreferredDiscountEligibilityRule) Execute(memory *engine.WorkingMemory) error {
	memory.AddDerivedFact(engine.DerivedFact{
		Name:  engine.PreferredDiscountEligibleFact,
		Value: memory.Quote.ID,
	})
	return nil
}
