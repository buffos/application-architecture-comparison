package readmodel

import (
	"testing"

	"rules-engine-architecture/internal/engine"
)

func TestProjectProductsSortsAndCopiesProductFacts(t *testing.T) {
	products := []engine.ProductFact{
		{
			ID:                "PRD-002",
			Name:              "Configured Workstation",
			Category:          "CustomBuild",
			UnitPriceCents:    125000,
			AvailableQuantity: 2,
			ShortagePolicy:    engine.StockShortageBackorder,
		},
		{
			ID:                "PRD-001",
			Name:              "Standard Workstation",
			Category:          "Standard",
			AvailableQuantity: 12,
		},
	}

	views := ProjectProducts(products)

	if len(views) != 2 || views[0].ID != "PRD-001" || views[1].ID != "PRD-002" {
		t.Fatalf("expected sorted product views, got %+v", views)
	}
	if views[1].AvailableQuantity != 2 || views[1].ShortagePolicy != engine.StockShortageBackorder {
		t.Fatalf("expected inventory data in view, got %+v", views[1])
	}

	products[1].AvailableQuantity = 0
	if views[1].AvailableQuantity != 2 {
		t.Fatal("expected product view to copy source values")
	}
}
