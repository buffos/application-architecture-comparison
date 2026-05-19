package application

import (
	"sort"

	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ReturnRateByCategoryRow struct {
	Category        string
	ShippedQuantity int
	ReturnQuantity  int
	ReturnRate      float64
}

type GetReturnRateByCategoryReportUseCase struct {
	orders  ports.OrderRepository
	returns ports.ReturnRequestRepository
}

func NewGetReturnRateByCategoryReportUseCase(orders ports.OrderRepository, returns ports.ReturnRequestRepository) GetReturnRateByCategoryReportUseCase {
	return GetReturnRateByCategoryReportUseCase{
		orders:  orders,
		returns: returns,
	}
}

func (uc GetReturnRateByCategoryReportUseCase) Execute() ([]ReturnRateByCategoryRow, error) {
	orders, err := uc.orders.ListByStatus(domain.OrderStatusShipped)
	if err != nil {
		return nil, err
	}

	requests, err := uc.returns.ListByStatus("")
	if err != nil {
		return nil, err
	}

	shippedByCategory := make(map[string]int)
	for _, order := range orders {
		for _, line := range order.Lines {
			shippedByCategory[line.ProductCategory] += line.Quantity
		}
	}

	returnedByCategory := make(map[string]int)
	for _, request := range requests {
		if request.Status != domain.ReturnStatusAccepted && request.Status != domain.ReturnStatusRefunded {
			continue
		}

		for _, line := range request.Lines {
			returnedByCategory[line.ProductCategory] += line.Quantity
		}
	}

	categories := make([]string, 0, len(shippedByCategory))
	for category := range shippedByCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)

	rows := make([]ReturnRateByCategoryRow, 0, len(categories))
	for _, category := range categories {
		row := ReturnRateByCategoryRow{
			Category:        category,
			ShippedQuantity: shippedByCategory[category],
			ReturnQuantity:  returnedByCategory[category],
		}
		if row.ShippedQuantity > 0 {
			row.ReturnRate = float64(row.ReturnQuantity) / float64(row.ShippedQuantity)
		}
		rows = append(rows, row)
	}

	return rows, nil
}
