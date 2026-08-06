package reports

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
	domainreturns "domain-driven-design-architecture/internal/domain/returns"
)

func TestBuildReturnRateByCategoryReport(t *testing.T) {
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, price)
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.Submit(); err != nil {
		t.Fatal(err)
	}
	order, err := ordering.NewOrderFromQuote("order-001", quote)
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if err := order.MarkShipped(); err != nil {
		t.Fatal(err)
	}
	request, err := domainreturns.NewReturnRequestFromShippedOrder("return-001", order, "damaged")
	if err != nil {
		t.Fatal(err)
	}
	if err := request.Accept(); err != nil {
		t.Fatal(err)
	}
	report := BuildReturnRateByCategoryReport([]ordering.Order{order}, []domainreturns.ReturnRequest{request})
	if len(report.Rows) != 1 || report.Rows[0].ShippedQuantity != 2 || report.Rows[0].ReturnedQuantity != 2 || report.Rows[0].ReturnRate != 1 {
		t.Fatalf("unexpected report %+v", report)
	}
}
