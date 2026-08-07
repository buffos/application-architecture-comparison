package records

import "sort"

// ReturnRateCategoryRow is one category's read-time return metric.
type ReturnRateCategoryRow struct {
	Category         string
	ShippedQuantity  int
	ReturnedQuantity int
	ReturnRate       float64
}

// ReturnRateByCategoryReport is a sorted category projection.
type ReturnRateByCategoryReport struct {
	Rows []ReturnRateCategoryRow
}

// GetReturnRateByCategoryReport scans order snapshots without changing them.
func GetReturnRateByCategoryReport(db *Database) (ReturnRateByCategoryReport, error) {
	if db == nil {
		return ReturnRateByCategoryReport{}, ErrDatabaseRequired
	}

	rowsByCategory := make(map[string]ReturnRateCategoryRow)
	for _, order := range db.orders {
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
		if row.ShippedQuantity <= 0 {
			continue
		}
		row.ReturnRate = float64(row.ReturnedQuantity) / float64(row.ShippedQuantity)
		rows = append(rows, row)
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].Category < rows[j].Category
	})
	return ReturnRateByCategoryReport{Rows: rows}, nil
}
