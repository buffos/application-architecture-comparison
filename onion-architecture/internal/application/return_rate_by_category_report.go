package application

import "onion-architecture/internal/domain"

type ReturnRateByCategoryRow struct {
	Category        string
	ShippedQuantity int
	ReturnedQuantity int
	ReturnRate      float64
}

type ReturnRateByCategoryReportService struct {
	orders  OrderFinder
	returns ReturnRequestFinder
}

func NewReturnRateByCategoryReportService(orders OrderFinder, returns ReturnRequestFinder) ReturnRateByCategoryReportService {
	return ReturnRateByCategoryReportService{
		orders:  orders,
		returns: returns,
	}
}

func (s ReturnRateByCategoryReportService) Execute() ([]ReturnRateByCategoryRow, error) {
	shippedOrders, err := s.orders.ListByStatus(domain.OrderStatusShipped)
	if err != nil {
		return nil, err
	}

	refundedReturns, err := s.returns.ListByStatus(domain.ReturnRequestStatusRefunded)
	if err != nil {
		return nil, err
	}

	type totals struct {
		shipped  int
		returned int
	}

	byCategory := make(map[string]*totals)
	ordersByID := make(map[string]domain.Order)

	for _, order := range shippedOrders {
		ordersByID[order.ID] = order
		for _, line := range order.Lines {
			entry := byCategory[line.ProductCategory]
			if entry == nil {
				entry = &totals{}
				byCategory[line.ProductCategory] = entry
			}
			entry.shipped += line.Quantity
		}
	}

	for _, request := range refundedReturns {
		order, ok := ordersByID[request.OrderID]
		if !ok {
			continue
		}

		for _, line := range order.Lines {
			entry := byCategory[line.ProductCategory]
			if entry == nil {
				entry = &totals{}
				byCategory[line.ProductCategory] = entry
			}
			entry.returned += line.Quantity
		}
	}

	result := make([]ReturnRateByCategoryRow, 0, len(byCategory))
	for category, totals := range byCategory {
		rate := 0.0
		if totals.shipped > 0 {
			rate = float64(totals.returned) / float64(totals.shipped)
		}

		result = append(result, ReturnRateByCategoryRow{
			Category:         category,
			ShippedQuantity:  totals.shipped,
			ReturnedQuantity: totals.returned,
			ReturnRate:       rate,
		})
	}

	return result, nil
}
