package reports

import (
	"testing"

	"rich-domain-model-architecture/internal/domain/ordering"
	domainreturns "rich-domain-model-architecture/internal/domain/returns"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestReturnRateReportCountsOnlyAcceptedReturns(t *testing.T) {
	order := shippedOrderForReport(t)
	accepted := returnForReport(t, domainreturns.ReviewDecisionAccept, "return-accepted")
	rejected := returnForReport(t, domainreturns.ReviewDecisionReject, "return-rejected")
	report := BuildReturnRateByCategoryReport([]ordering.Order{order}, []domainreturns.ReturnRequest{accepted, rejected})
	if len(report.Rows) != 1 {
		t.Fatalf("rows = %+v", report.Rows)
	}
	row := report.Rows[0]
	if row.Category != "Standard" || row.ShippedQuantity != 2 || row.ReturnedQuantity != 2 || row.ReturnRate != 1 {
		t.Fatalf("row = %+v", row)
	}
}

func TestReturnRateReportHandlesNoOrders(t *testing.T) {
	report := BuildReturnRateByCategoryReport(nil, nil)
	if len(report.Rows) != 0 {
		t.Fatalf("rows = %+v", report.Rows)
	}
}

func shippedOrderForReport(t *testing.T) ordering.Order {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLineWithCategory("sku-001", quoting.ProductCategoryStandard, 2, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-report", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	order, err := ordering.NewOrderFromQuote("order-report", quote)
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if err := order.MarkShipped(); err != nil {
		t.Fatal(err)
	}
	return order
}

func returnForReport(t *testing.T, decision domainreturns.ReviewDecision, id domainreturns.ReturnRequestID) domainreturns.ReturnRequest {
	t.Helper()
	request, err := domainreturns.NewReturnRequestFromShippedOrder(id, shippedOrderForReport(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	if err := request.Review(decision); err != nil {
		t.Fatal(err)
	}
	return request
}
