package payments

import "testing"

type stubGateway struct {
	request PaymentRequest
	err     error
}

func (g *stubGateway) Capture(request PaymentRequest) error {
	if g.err != nil {
		return g.err
	}

	g.request = request
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

	if gateway.request.OrderID != "order-001" {
		t.Fatalf("expected order-001, got %s", gateway.request.OrderID)
	}
}
