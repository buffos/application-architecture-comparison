package reporting

import (
	"sort"

	"rules-engine-architecture/internal/engine"
)

type LowStockRow struct {
	ProductID         string
	ProductCategory   string
	AvailableQuantity int
	Threshold         int
}

func BuildLowStockReport(products []engine.ProductFact, threshold int) []LowStockRow {
	if threshold <= 0 {
		return []LowStockRow{}
	}

	rows := make([]LowStockRow, 0)
	for _, product := range products {
		if product.AvailableQuantity >= threshold {
			continue
		}
		rows = append(rows, LowStockRow{
			ProductID:         product.ID,
			ProductCategory:   product.Category,
			AvailableQuantity: product.AvailableQuantity,
			Threshold:         threshold,
		})
	}

	sort.Slice(rows, func(left, right int) bool {
		return rows[left].ProductID < rows[right].ProductID
	})
	return rows
}
