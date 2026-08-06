package returns

import (
	"errors"
	"testing"

	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
	domainreturns "domain-driven-design-architecture/internal/domain/returns"
)

func TestReaderProjectsReturnAggregate(t *testing.T) {
	request, err := domainreturns.NewReturnRequestFromShippedOrder("return-001", shippedOrderForQuery(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	if err := request.AssignRequester("agent-001"); err != nil {
		t.Fatal(err)
	}
	if err := request.ReviewBy(domainreturns.ReviewDecisionAccept, "reviewer-001"); err != nil {
		t.Fatal(err)
	}
	if err := request.ProcessBy("processor-001"); err != nil {
		t.Fatal(err)
	}
	reader := NewInMemoryReader()
	reader.Save(request)
	details, err := reader.GetReturnRequest("return-001")
	if err != nil {
		t.Fatal(err)
	}
	if details.Status != string(domainreturns.ReturnStatusAccepted) || details.ReviewedBy != "reviewer-001" || len(details.Lines) != 1 {
		t.Fatalf("unexpected details %+v", details)
	}
	if got := reader.ListReturnRequests(domainreturns.ReturnStatusAccepted); len(got) != 1 || got[0].ReturnRequestID != "return-001" {
		t.Fatalf("unexpected summaries %+v", got)
	}
	if _, err := reader.GetReturnRequest("missing"); !errors.Is(err, ErrReturnRequestNotFound) {
		t.Fatalf("missing request returned %v", err)
	}
}

func shippedOrderForQuery(t *testing.T) ordering.Order {
	t.Helper()
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 1, price)
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
