package returns

import (
	"errors"
	"testing"

	domainreturns "rich-domain-model-architecture/internal/domain/returns"
	"rich-domain-model-architecture/internal/domain/ordering"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestReturnQueryProjectsDetailsAndFiltersSummaries(t *testing.T) {
	request := queryReturnRequest(t)
	reader := NewInMemoryReader()
	reader.Save(request)

	details, err := reader.GetReturnRequest(string(request.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if details.Status != string(domainreturns.ReturnStatusAccepted) || details.ReviewedBy != "reviewer-001" || len(details.Lines) != 1 {
		t.Fatalf("details = %+v", details)
	}
	rows := reader.ListReturnRequests(domainreturns.ReturnStatusAccepted)
	if len(rows) != 1 || rows[0].ReturnRequestID != string(request.ID()) {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestReturnQueryReturnsNotFound(t *testing.T) {
	if _, err := NewInMemoryReader().GetReturnRequest("missing"); !errors.Is(err, ErrReturnRequestNotFound) {
		t.Fatalf("missing query returned %v", err)
	}
}

func queryReturnRequest(t *testing.T) domainreturns.ReturnRequest {
	t.Helper()
	request, err := domainreturns.NewReturnRequestFromShippedOrder("return-query", shippedOrderForQuery(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	if err := request.AssignRequester("customer-001"); err != nil {
		t.Fatal(err)
	}
	if err := request.ReviewBy(domainreturns.ReviewDecisionAccept, "reviewer-001"); err != nil {
		t.Fatal(err)
	}
	if err := request.ProcessBy("processor-001"); err != nil {
		t.Fatal(err)
	}
	return request
}

func shippedOrderForQuery(t *testing.T) ordering.Order {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-query", "customer-001")
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
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if err := order.MarkShipped(); err != nil {
		t.Fatal(err)
	}
	return order
}
