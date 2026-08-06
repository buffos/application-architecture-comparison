package reporting

import (
	"component-based-architecture/internal/components/orders"
	"component-based-architecture/internal/components/returns"
	"testing"
)

type returnRateOrderReader struct{}

func (returnRateOrderReader) ListOrders(orders.ListOrdersQuery) []orders.OrderSummary {
	return []orders.OrderSummary{{OrderID: "order-001", Status: orders.OrderStatusShipped}}
}
func (returnRateOrderReader) GetOrder(orders.GetOrderQuery) (orders.OrderDetails, error) {
	return orders.OrderDetails{Lines: []orders.OrderLineDetails{{ProductCategory: "Standard", Quantity: 4}}}, nil
}

type returnRateReturnReader struct{}

func (returnRateReturnReader) ListReturnRequests(returns.ListReturnRequestsQuery) []returns.ReturnRequestSummary {
	return []returns.ReturnRequestSummary{{ReturnRequestID: "return-001", Status: returns.ReturnRequestStatusRefunded}}
}
func (returnRateReturnReader) GetReturnRequest(returns.GetReturnRequestQuery) (returns.ReturnRequestDetails, error) {
	return returns.ReturnRequestDetails{Lines: []returns.ReturnLineDetails{{ProductCategory: "Standard", Quantity: 1}}}, nil
}

func TestReturnRateByCategoryReportAggregatesLines(t *testing.T) {
	report, err := NewReturnRateComponent(returnRateOrderReader{}, returnRateReturnReader{}).ReturnRateByCategoryReport()
	if err != nil || len(report.Rows) != 1 || report.Rows[0].ShippedQuantity != 4 || report.Rows[0].ReturnedQuantity != 1 || report.Rows[0].ReturnRate != 0.25 {
		t.Fatalf("report=%+v err=%v", report, err)
	}
}
