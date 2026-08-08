package reporting

import "sort"

type ReturnRecord struct {
	ProductCategory string
	Accepted        bool
}

type ReturnRateRow struct {
	ProductCategory   string
	AttemptedReturns  int
	AcceptedReturns   int
	ReturnRatePercent float64
}

func BuildReturnRateByCategory(records []ReturnRecord) []ReturnRateRow {
	grouped := map[string]*ReturnRateRow{}
	for _, record := range records {
		category := record.ProductCategory
		if category == "" {
			category = "Unknown"
		}
		row, found := grouped[category]
		if !found {
			row = &ReturnRateRow{ProductCategory: category}
			grouped[category] = row
		}

		row.AttemptedReturns++
		if record.Accepted {
			row.AcceptedReturns++
		}
	}

	rows := make([]ReturnRateRow, 0, len(grouped))
	for _, row := range grouped {
		row.ReturnRatePercent = float64(row.AcceptedReturns) * 100 / float64(row.AttemptedReturns)
		rows = append(rows, *row)
	}
	sort.Slice(rows, func(left, right int) bool {
		return rows[left].ProductCategory < rows[right].ProductCategory
	})
	return rows
}
