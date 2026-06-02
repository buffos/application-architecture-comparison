package payments

import "testing"

type stubGateway struct {
	captureRequest PaymentRequest
	refundRequest  RefundRequest
	err            error
}

func (g *stubGateway) Capture(request PaymentRequest) error {
	if g.err != nil {
		return g.err
	}

	g.captureRequest = request
	return nil
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

	err := service.Capture(PaymentRequest{
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
