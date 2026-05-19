package application

import (
	"sort"

	"hexagonal-architecture/internal/core/ports"
)

type TopDiscountedProductRow struct {
	SKU                 string
	ProductName         string
	QuotedQuantity      int
	TotalDiscountAmount int
	AverageDiscountRate float64
}

type GetTopDiscountedProductsReportUseCase struct {
	quotes ports.QuoteRepository
}

func NewGetTopDiscountedProductsReportUseCase(quotes ports.QuoteRepository) GetTopDiscountedProductsReportUseCase {
	return GetTopDiscountedProductsReportUseCase{quotes: quotes}
}

func (uc GetTopDiscountedProductsReportUseCase) Execute() ([]TopDiscountedProductRow, error) {
	quotes, err := uc.quotes.ListByStatus("")
	if err != nil {
		return nil, err
	}

	type aggregate struct {
		sku            string
		productName    string
		quotedQuantity int
		baseAmount     int
		discountAmount int
	}

	bySKU := make(map[string]*aggregate)
	for _, quote := range quotes {
		for _, line := range quote.Lines {
			row, ok := bySKU[line.SKU]
			if !ok {
				row = &aggregate{
					sku:         line.SKU,
					productName: line.ProductName,
				}
				bySKU[line.SKU] = row
			}

			baseLineAmount := line.BaseUnitPrice * line.Quantity
			discountLineAmount := (line.BaseUnitPrice - line.AdjustedUnitPrice) * line.Quantity

			row.quotedQuantity += line.Quantity
			row.baseAmount += baseLineAmount
			row.discountAmount += discountLineAmount
		}
	}

	rows := make([]TopDiscountedProductRow, 0, len(bySKU))
	for _, row := range bySKU {
		reportRow := TopDiscountedProductRow{
			SKU:                 row.sku,
			ProductName:         row.productName,
			QuotedQuantity:      row.quotedQuantity,
			TotalDiscountAmount: row.discountAmount,
		}
		if row.baseAmount > 0 {
			reportRow.AverageDiscountRate = float64(row.discountAmount) / float64(row.baseAmount)
		}
		rows = append(rows, reportRow)
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].TotalDiscountAmount == rows[j].TotalDiscountAmount {
			return rows[i].SKU < rows[j].SKU
		}
		return rows[i].TotalDiscountAmount > rows[j].TotalDiscountAmount
	})

	return rows, nil
}
