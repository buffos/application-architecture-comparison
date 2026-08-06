package reporting

import (
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/returns"
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

type ReturnRateComponent struct {
	orders  orders.Reader
	returns returns.Reader
}

func NewReturnRateComponent(orders orders.Reader, returns returns.Reader) *ReturnRateComponent {
	return &ReturnRateComponent{orders: orders, returns: returns}
}

func (c *ReturnRateComponent) ReturnRateByCategoryReport() (ReturnRateByCategoryReport, error) {
	rows := make(map[string]*ReturnRateByCategoryRow)
	for _, summary := range c.orders.ListOrders(orders.ListOrdersQuery{Status: orders.OrderStatusShipped}) {
		details, err := c.orders.GetOrder(orders.GetOrderQuery{OrderID: summary.OrderID})
		if err != nil {
			return ReturnRateByCategoryReport{}, err
		}
		for _, line := range details.Lines {
			row := rows[line.ProductCategory]
			if row == nil {
				row = &ReturnRateByCategoryRow{Category: line.ProductCategory}
				rows[line.ProductCategory] = row
			}
			row.ShippedQuantity += line.Quantity
		}
	}
	for _, summary := range c.returns.ListReturnRequests(returns.ListReturnRequestsQuery{Status: returns.ReturnRequestStatusRefunded}) {
		details, err := c.returns.GetReturnRequest(returns.GetReturnRequestQuery{ReturnRequestID: summary.ReturnRequestID})
		if err != nil {
			return ReturnRateByCategoryReport{}, err
		}
		for _, line := range details.Lines {
			row := rows[line.ProductCategory]
			if row == nil {
				row = &ReturnRateByCategoryRow{Category: line.ProductCategory}
				rows[line.ProductCategory] = row
			}
			row.ReturnedQuantity += line.Quantity
		}
	}
	report := ReturnRateByCategoryReport{Rows: make([]ReturnRateByCategoryRow, 0, len(rows))}
	for _, row := range rows {
		if row.ShippedQuantity > 0 {
			row.ReturnRate = float64(row.ReturnedQuantity) / float64(row.ShippedQuantity)
		}
		report.Rows = append(report.Rows, *row)
	}
	return report, nil
}
