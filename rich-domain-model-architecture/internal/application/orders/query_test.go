package orders

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/ordering"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestOrderQueryProjectsDetailsAndFiltersByStatus(t *testing.T) {
	order := queryOrder(t)
	reader := NewInMemoryReader()
	if err := reader.Save(order); err != nil {
		t.Fatal(err)
	}
	details, err := reader.GetOrder(string(order.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if details.TotalCents != 2000 || details.Currency != "USD" || len(details.Lines) != 1 {
		t.Fatalf("details = %+v", details)
	}
	rows := reader.ListOrders(order.Status())
	if len(rows) != 1 || rows[0].OrderID != string(order.ID()) {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestOrderQueryReturnsNotFound(t *testing.T) {
	if _, err := NewInMemoryReader().GetOrder("missing"); !errors.Is(err, ErrOrderNotFound) {
		t.Fatalf("missing query returned %v", err)
	}
}

func queryOrder(t *testing.T) ordering.Order {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-order-query", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	order, err := ordering.NewOrderFromQuote("order-query", quote)
	if err != nil {
		t.Fatal(err)
	}
	return order
}
