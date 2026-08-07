package returns

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/ordering"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func shippedOrderForReturn(t *testing.T) ordering.Order {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLineWithCategory("sku-001", quoting.ProductCategoryStandard, 2, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-001", "customer-001")
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

func TestReturnRequestRequiresShippedOrderAndCopiesLines(t *testing.T) {
	pendingOrder := pendingOrderForReturn(t)
	if _, err := NewReturnRequestFromShippedOrder("return-001", pendingOrder, "damaged"); !errors.Is(err, ErrOrderNotShipped) {
		t.Fatalf("pending order returned %v", err)
	}
	request, err := NewReturnRequestFromShippedOrder("return-001", shippedOrderForReturn(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	if request.Status() != ReturnStatusRequested || len(request.Lines()) != 1 || request.Lines()[0].Quantity() != 2 {
		t.Fatalf("unexpected return request: status=%s lines=%+v", request.Status(), request.Lines())
	}
}

func pendingOrderForReturn(t *testing.T) ordering.Order {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-pending", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	order, err := ordering.NewOrderFromQuote("order-pending", quote)
	if err != nil {
		t.Fatal(err)
	}
	return order
}

func TestRefundCanBeIssuedOnlyOnce(t *testing.T) {
	amount, err := NewMoney(2000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	refund, err := NewRefund("refund-001", "return-001", amount)
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
