package returns

import (
	"errors"
	"testing"

	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
	domainreturns "domain-driven-design-architecture/internal/domain/returns"
)

func TestReviewServiceReplaysCompletedResult(t *testing.T) {
	request, err := domainreturns.NewReturnRequestFromShippedOrder("return-001", shippedOrder(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	service := NewReviewService(NewInMemoryIdempotencyStore())
	first, err := service.Review(&request, domainreturns.ReviewDecisionAccept, "reviewer-001", "processor-001", "review-001")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Review(&request, domainreturns.ReviewDecisionAccept, "reviewer-002", "processor-002", "review-001")
	if err != nil {
		t.Fatal(err)
	}
	if first != second || request.ReviewedBy() != "reviewer-001" || request.ProcessedBy() != "processor-001" {
		t.Fatalf("retry changed result: first=%+v second=%+v request=%+v", first, second, request)
	}
}

func TestReviewServiceRequiresKey(t *testing.T) {
	request, err := domainreturns.NewReturnRequestFromShippedOrder("return-001", shippedOrder(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	service := NewReviewService(NewInMemoryIdempotencyStore())
	if _, err := service.Review(&request, domainreturns.ReviewDecisionReject, "reviewer-001", "", ""); !errors.Is(err, ErrIdempotencyKeyRequired) {
		t.Fatalf("missing key returned %v", err)
	}
}

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
