package orders

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
)

func TestReaderProjectsOrderAggregate(t *testing.T) {
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
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	order, err := ordering.NewOrderFromQuote("order-001", quote)
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	reader := NewInMemoryReader()
	if err := reader.Save(order); err != nil {
		t.Fatal(err)
	}
	details, err := reader.GetOrder("order-001")
	if err != nil {
		t.Fatal(err)
	}
	if details.Status != string(ordering.OrderStatusPaid) || details.TotalCents != 2000 || len(details.Lines) != 1 {
		t.Fatalf("unexpected details %+v", details)
	}
	if got := reader.ListOrders(ordering.OrderStatusPaid); len(got) != 1 || got[0].OrderID != "order-001" {
		t.Fatalf("unexpected summaries %+v", got)
	}
}
