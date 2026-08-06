package returns

import (
	"errors"
	"testing"

	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
)

func shippedOrder(t *testing.T) ordering.Order {
	t.Helper()
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
	if err := order.MarkShipped(); err != nil {
		t.Fatal(err)
	}
	return order
}

func TestReturnRequestAndRefundLifecycles(t *testing.T) {
	request, err := NewReturnRequestFromShippedOrder("return-001", shippedOrder(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	if request.Status() != ReturnStatusRequested || len(request.Lines()) != 1 || request.Lines()[0].Quantity() != 2 {
		t.Fatalf("unexpected request %+v", request)
	}
	if err := request.Accept(); err != nil {
		t.Fatal(err)
	}
	amount, err := NewMoney(2000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	refund, err := NewRefund("refund-001", request.ID(), amount)
	if err != nil {
		t.Fatal(err)
	}
	if err := refund.Issue(); err != nil {
		t.Fatal(err)
	}
	if refund.Status() != RefundStatusIssued {
		t.Fatalf("status = %s, want %s", refund.Status(), RefundStatusIssued)
	}
	if err := refund.Issue(); !errors.Is(err, ErrRefundNotIssuable) {
		t.Fatalf("repeated issue returned %v", err)
	}
}

func TestReturnRequestRequiresShippedOrder(t *testing.T) {
	if _, err := NewReturnRequestFromShippedOrder("return-001", ordering.Order{}, "damaged"); !errors.Is(err, ErrOrderNotShipped) {
		t.Fatalf("unshipped return returned %v", err)
	}
}
