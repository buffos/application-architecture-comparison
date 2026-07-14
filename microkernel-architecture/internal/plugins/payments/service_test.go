package payments

import "testing"

func TestCapture(t *testing.T) {
	service := NewService()

	if err := service.Capture("order-001", 30000); err != nil {
		t.Fatalf("expected capture to succeed, got %v", err)
	}
}

func TestRefund(t *testing.T) {
	service := NewService()

	if err := service.Refund("order-001", 30000); err != nil {
		t.Fatalf("expected refund to succeed, got %v", err)
	}
}
