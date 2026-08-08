package rules

import (
	"fmt"

	"rules-engine-architecture/internal/engine"
)

// InventoryShortageRule turns a stock shortage into the policy selected by the
// ProductFact. It does not reserve or mutate inventory.
type InventoryShortageRule struct{}

func (InventoryShortageRule) Name() string {
	return "Inventory Shortage Rule"
}

func (InventoryShortageRule) Priority() int {
	return 80
}

func (InventoryShortageRule) ConflictGroup() string {
	return ""
}

func (InventoryShortageRule) Evaluate(memory *engine.WorkingMemory) bool {
	for _, line := range memory.Quote.Lines {
		product, found := productForLine(memory, line.ProductID)
		if found && line.Quantity > product.AvailableQuantity {
			return true
		}
	}
	return false
}

func (rule InventoryShortageRule) Execute(memory *engine.WorkingMemory) error {
	for _, line := range memory.Quote.Lines {
		product, found := productForLine(memory, line.ProductID)
		if !found || line.Quantity <= product.AvailableQuantity {
			continue
		}

		finding := engine.Finding{
			RuleName: rule.Name(),
			Message: fmt.Sprintf(
				"Product %s requests %d units but only %d are available",
				product.ID,
				line.Quantity,
				product.AvailableQuantity,
			),
		}

		switch product.ShortagePolicy {
		case engine.StockShortageBackorder:
			finding.Severity = "inventory-backorder"
		case engine.StockShortageReject:
			finding.Severity = "inventory-rejected"
		default:
			return fmt.Errorf(
				"product %s has unsupported stock shortage policy %q",
				product.ID,
				product.ShortagePolicy,
			)
		}

		memory.AddFinding(finding)
	}

	return nil
}

func productForLine(memory *engine.WorkingMemory, productID string) (engine.ProductFact, bool) {
	for _, product := range memory.Products {
		if product.ID == productID {
			return product, true
		}
	}
	return engine.ProductFact{}, false
}
