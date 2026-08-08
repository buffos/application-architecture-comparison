package readmodel

import (
	"sort"

	"rules-engine-architecture/internal/engine"
)

type ProductView struct {
	ID                string
	Name              string
	Category          string
	UnitPriceCents    int64
	AvailableQuantity int
	ShortagePolicy    engine.StockShortagePolicy
}

func ProjectProducts(products []engine.ProductFact) []ProductView {
	views := make([]ProductView, 0, len(products))
	for _, product := range products {
		views = append(views, ProductView{
			ID:                product.ID,
			Name:              product.Name,
			Category:          product.Category,
			UnitPriceCents:    product.UnitPriceCents,
			AvailableQuantity: product.AvailableQuantity,
			ShortagePolicy:    product.ShortagePolicy,
		})
	}

	sort.Slice(views, func(left, right int) bool {
		return views[left].ID < views[right].ID
	})
	return views
}
