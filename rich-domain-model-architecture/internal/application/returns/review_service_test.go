package returns

import (
	"errors"
	"testing"

	domainreturns "rich-domain-model-architecture/internal/domain/returns"
	"rich-domain-model-architecture/internal/domain/ordering"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestReviewServiceReplaysSameResultForRepeatedKey(t *testing.T) {
	request := requestForReviewService(t)
	service := NewReviewService(NewInMemoryIdempotencyStore())
	if err := request.AssignRequester("customer-001"); err != nil {
		t.Fatal(err)
	}

	first, err := service.Review(&request, domainreturns.ReviewDecisionAccept, "reviewer-001", "processor-001", "review-key-001")
	if err != nil {
		t.Fatal(err)
	}
	second, err := service.Review(&request, domainreturns.ReviewDecisionReject, "reviewer-002", "processor-002", "review-key-001")
	if err != nil {
		t.Fatal(err)
	}
	if first != second {
		t.Fatalf("replayed result = %+v, first = %+v", second, first)
	}
	if request.Status() != domainreturns.ReturnStatusAccepted || request.ReviewedBy() != "reviewer-001" {
		t.Fatalf("retry changed aggregate: status=%s reviewer=%s", request.Status(), request.ReviewedBy())
	}
}

func TestReviewServiceRequiresIdempotencyKey(t *testing.T) {
	request := requestForReviewService(t)
	service := NewReviewService(NewInMemoryIdempotencyStore())
	if _, err := service.Review(&request, domainreturns.ReviewDecisionAccept, "reviewer-001", "processor-001", ""); !errors.Is(err, ErrIdempotencyKeyRequired) {
		t.Fatalf("missing key returned %v", err)
	}
}

func requestForReviewService(t *testing.T) domainreturns.ReturnRequest {
	t.Helper()
	request, err := domainreturns.NewReturnRequestFromShippedOrder("return-service", shippedOrderForReviewService(t), "damaged")
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func shippedOrderForReviewService(t *testing.T) ordering.Order {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-service", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	order, err := ordering.NewOrderFromQuote("order-service", quote)
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
