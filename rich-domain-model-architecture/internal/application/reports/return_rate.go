package reports

import (
	"sort"

	"rich-domain-model-architecture/internal/domain/ordering"
	domainreturns "rich-domain-model-architecture/internal/domain/returns"
)

type ReturnRateByCategoryRow struct {
	Category         string
	ShippedQuantity  int
	ReturnedQuantity int
	ReturnRate       float64
}

type ReturnRateByCategoryReport struct {
	Rows []ReturnRateByCategoryRow
}

func BuildReturnRateByCategoryReport(orders []ordering.Order, requests []domainreturns.ReturnRequest) ReturnRateByCategoryReport {
	rows := make(map[string]*ReturnRateByCategoryRow)
	for _, order := range orders {
		if order.Status() != ordering.OrderStatusShipped {
			continue
		}
		for _, line := range order.Lines() {
			category := string(line.ProductCategory())
			if rows[category] == nil {
				rows[category] = &ReturnRateByCategoryRow{Category: category}
			}
			rows[category].ShippedQuantity += line.Quantity()
		}
	}
	for _, request := range requests {
		if request.Status() != domainreturns.ReturnStatusAccepted {
			continue
		}
		for _, line := range request.Lines() {
			category := string(line.ProductCategory())
			if rows[category] == nil {
				rows[category] = &ReturnRateByCategoryRow{Category: category}
			}
			rows[category].ReturnedQuantity += line.Quantity()
		}
	}
	result := ReturnRateByCategoryReport{Rows: make([]ReturnRateByCategoryRow, 0, len(rows))}
	for _, row := range rows {
		if row.ShippedQuantity > 0 {
			row.ReturnRate = float64(row.ReturnedQuantity) / float64(row.ShippedQuantity)
		}
		result.Rows = append(result.Rows, *row)
	}
	sort.Slice(result.Rows, func(i, j int) bool { return result.Rows[i].Category < result.Rows[j].Category })
	return result
}
