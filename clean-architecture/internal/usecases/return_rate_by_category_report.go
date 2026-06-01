package usecases

import "clean-architecture/internal/entities"

type ReturnRateByCategoryReportInput struct{}

type ReturnRateByCategoryItem struct {
	Category         string
	ShippedQuantity  int
	ReturnedQuantity int
	ReturnRate       float64
}

type ReturnRateByCategoryReportOutput struct {
	Categories []ReturnRateByCategoryItem
}

type ReturnRateByCategoryReportInputBoundary interface {
	Execute(input ReturnRateByCategoryReportInput) error
}

type ReturnRateByCategoryReportOutputBoundary interface {
	Present(output ReturnRateByCategoryReportOutput) error
}

type ReturnReportOrderReader interface {
	ListByStatus(status string) ([]entities.Order, error)
}

type ReturnReportRequestReader interface {
	ListByStatus(status string) ([]entities.ReturnRequest, error)
}

type ReturnReportProductReader interface {
	FindBySKU(sku string) (entities.Product, error)
}

type ReturnRateByCategoryReportInteractor struct {
	orders   ReturnReportOrderReader
	returns  ReturnReportRequestReader
	products ReturnReportProductReader
	output   ReturnRateByCategoryReportOutputBoundary
}

func NewReturnRateByCategoryReportInteractor(orders ReturnReportOrderReader, returns ReturnReportRequestReader, products ReturnReportProductReader, output ReturnRateByCategoryReportOutputBoundary) ReturnRateByCategoryReportInteractor {
	return ReturnRateByCategoryReportInteractor{
		orders:   orders,
		returns:  returns,
		products: products,
		output:   output,
	}
}

func (uc ReturnRateByCategoryReportInteractor) Execute(input ReturnRateByCategoryReportInput) error {
	_ = input

	shippedOrders, err := uc.orders.ListByStatus(entities.OrderStatusShipped)
	if err != nil {
		return err
	}

	refundedReturns, err := uc.returns.ListByStatus(entities.ReturnRequestStatusRefunded)
	if err != nil {
		return err
	}

	byCategory := make(map[string]*ReturnRateByCategoryItem)
	ordersByID := make(map[string]entities.Order, len(shippedOrders))

	for _, order := range shippedOrders {
		ordersByID[order.ID] = order

		for _, line := range order.Lines {
			product, err := uc.products.FindBySKU(line.SKU)
			if err != nil {
				return err
			}

			item := ensureCategory(byCategory, product.Category)
			item.ShippedQuantity += line.Quantity
		}
	}

	for _, request := range refundedReturns {
		order, ok := ordersByID[request.OrderID]
		if !ok {
			continue
		}

		for _, line := range order.Lines {
			product, err := uc.products.FindBySKU(line.SKU)
			if err != nil {
				return err
			}

			item := ensureCategory(byCategory, product.Category)
			item.ReturnedQuantity += line.Quantity
		}
	}

	categories := make([]ReturnRateByCategoryItem, 0, len(byCategory))
	for _, item := range byCategory {
		if item.ShippedQuantity > 0 {
			item.ReturnRate = float64(item.ReturnedQuantity) / float64(item.ShippedQuantity)
		}
		categories = append(categories, *item)
	}

	for i := 0; i < len(categories)-1; i++ {
		for j := i + 1; j < len(categories); j++ {
			if categories[j].Category < categories[i].Category {
				categories[i], categories[j] = categories[j], categories[i]
			}
		}
	}

	return uc.output.Present(ReturnRateByCategoryReportOutput{
		Categories: categories,
	})
}

func ensureCategory(byCategory map[string]*ReturnRateByCategoryItem, category string) *ReturnRateByCategoryItem {
	item, ok := byCategory[category]
	if !ok {
		item = &ReturnRateByCategoryItem{Category: category}
		byCategory[category] = item
	}

	return item
}
