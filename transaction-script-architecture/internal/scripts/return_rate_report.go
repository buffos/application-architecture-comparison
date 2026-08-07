package scripts

import (
	"sort"

	"transaction-script-architecture/internal/data"
)

type ReturnRateCategoryRow struct {
	Category         string
	ShippedQuantity  int
	ReturnedQuantity int
	ReturnRate       float64
}

type ReturnRateByCategoryReport struct {
	Rows []ReturnRateCategoryRow
}

func GetReturnRateByCategoryReport(store *data.Store) (ReturnRateByCategoryReport, error) {
	if store == nil {
		return ReturnRateByCategoryReport{}, ErrStoreRequired
	}

	rowsByCategory := make(map[string]ReturnRateCategoryRow)
	for _, order := range store.Orders {
		for _, line := range order.Lines {
			category := line.ProductCategory
			if category == "" {
				category = "Unknown"
			}

			row := rowsByCategory[category]
			row.Category = category
			row.ShippedQuantity += line.ShippedQuantity
			row.ReturnedQuantity += line.ReturnedQuantity
			rowsByCategory[category] = row
		}
	}

	rows := make([]ReturnRateCategoryRow, 0, len(rowsByCategory))
	for _, row := range rowsByCategory {
		if row.ShippedQuantity > 0 {
			row.ReturnRate = float64(row.ReturnedQuantity) / float64(row.ShippedQuantity)
		}
		rows = append(rows, row)
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Category < rows[j].Category
	})

	return ReturnRateByCategoryReport{Rows: rows}, nil
}
