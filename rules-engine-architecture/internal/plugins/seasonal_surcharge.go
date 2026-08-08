package plugins

import (
	"fmt"

	"rules-engine-architecture/internal/engine"
)

// SeasonalSurchargeRule is an example of a Rule supplied by an optional
// plugin package. It uses only the existing Rule contract.
type SeasonalSurchargeRule struct {
	category string
	percent  int
}

func NewSeasonalSurchargeRule(category string, percent int) SeasonalSurchargeRule {
	return SeasonalSurchargeRule{category: category, percent: percent}
}

func (SeasonalSurchargeRule) Name() string {
	return "Seasonal Surcharge Plugin Rule"
}

func (SeasonalSurchargeRule) Priority() int {
	return 45
}

func (SeasonalSurchargeRule) ConflictGroup() string {
	return ""
}

func (rule SeasonalSurchargeRule) Evaluate(memory *engine.WorkingMemory) bool {
	if rule.percent <= 0 || rule.category == "" {
		return false
	}

	for _, line := range memory.Quote.Lines {
		product, found := productForPlugin(memory, line.ProductID)
		if found && product.Category == rule.category {
			return true
		}
	}
	return false
}

func (rule SeasonalSurchargeRule) Execute(memory *engine.WorkingMemory) error {
	var matchingSubtotal int64
	for _, line := range memory.Quote.Lines {
		product, found := productForPlugin(memory, line.ProductID)
		if !found || product.Category != rule.category {
			continue
		}
		matchingSubtotal += int64(line.Quantity) * line.UnitPriceCents
	}

	adjustment := matchingSubtotal * int64(rule.percent) / 100
	memory.AddFinding(engine.Finding{
		RuleName:        rule.Name(),
		Severity:        "pricing-surcharge",
		AdjustmentCents: adjustment,
		Message: fmt.Sprintf(
			"Seasonal surcharge of %d%% adds %s to %s products",
			rule.percent,
			formatCents(adjustment),
			rule.category,
		),
	})
	return nil
}

func productForPlugin(memory *engine.WorkingMemory, productID string) (engine.ProductFact, bool) {
	for _, product := range memory.Products {
		if product.ID == productID {
			return product, true
		}
	}
	return engine.ProductFact{}, false
}

func formatCents(amountCents int64) string {
	return fmt.Sprintf("%d.%02d", amountCents/100, amountCents%100)
}
