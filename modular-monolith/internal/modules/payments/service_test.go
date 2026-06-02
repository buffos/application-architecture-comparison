package payments

import "testing"

type stubGateway struct {
	captureRequest PaymentRequest
	captureResult  CaptureResult
	refundRequest  RefundRequest
	err            error
}

func (g *stubGateway) Capture(request PaymentRequest) (CaptureResult, error) {
	if g.err != nil {
		return CaptureResult{}, g.err
	}

	g.captureRequest = request
	if g.captureResult.Outcome == "" {
		g.captureResult = CaptureResult{Outcome: CaptureOutcomeApproved}
	}
	return g.captureResult, nil
}

func (g *stubGateway) Refund(request RefundRequest) error {
	if g.err != nil {
		return g.err
	}

	g.refundRequest = request
	return nil
}

func TestCaptureDelegatesToGateway(t *testing.T) {
	gateway := &stubGateway{}
	service := NewService(gateway)

	result, err := service.Capture(PaymentRequest{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Amount:     30000,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gateway.captureRequest.OrderID != "order-001" {
		t.Fatalf("expected order-001, got %s", gateway.captureRequest.OrderID)
	}

	if result.Outcome != CaptureOutcomeApproved {
		t.Fatalf("expected approved outcome, got %s", result.Outcome)
	}
}

func TestRefundDelegatesToGateway(t *testing.T) {
	gateway := &stubGateway{}
	service := NewService(gateway)

	err := service.Refund(RefundRequest{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Amount:     30000,
		Reason:     "damaged item",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if gateway.refundRequest.Reason != "damaged item" {
		t.Fatalf("expected damaged item, got %s", gateway.refundRequest.Reason)
	}
}
