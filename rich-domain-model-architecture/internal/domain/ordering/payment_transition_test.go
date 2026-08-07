package ordering

import (
	"errors"
	"testing"
)

func TestOrderBecomesPaidOnlyFromPendingPayment(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusPaid {
		t.Fatalf("status = %s, want %s", order.Status(), OrderStatusPaid)
	}
	if err := order.MarkPaid(); !errors.Is(err, ErrOrderNotAwaitingPayment) {
		t.Fatalf("repeated paid transition returned %v", err)
	}
}

func TestOrderPaymentReviewBlocksShipmentUntilApproval(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaymentReview(); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusPaymentReview {
		t.Fatalf("status = %s, want %s", order.Status(), OrderStatusPaymentReview)
	}
	if err := order.MarkShipped(); !errors.Is(err, ErrOrderNotShippable) {
		t.Fatalf("shipped during review returned %v", err)
	}
	if err := order.ApprovePaymentReview(); err != nil {
		t.Fatal(err)
	}
	if err := order.MarkShipped(); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusShipped {
		t.Fatalf("status = %s, want %s", order.Status(), OrderStatusShipped)
	}
}

func TestOrderTracksPartialShipmentProgress(t *testing.T) {
	order, err := NewOrderFromQuote("order-001", approvedQuote(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	selection := []ShipmentSelection{{ProductSKU: "sku-001", Quantity: 1}}
	if err := order.ApplyShipment(selection); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusPartiallyShipped || order.Lines()[0].ShippedQuantity() != 1 {
		t.Fatalf("unexpected partial state: status=%s line=%+v", order.Status(), order.Lines()[0])
	}
	if err := order.Cancel(); !errors.Is(err, ErrOrderNotCancellable) {
		t.Fatalf("partial shipment cancellation returned %v", err)
	}
	if err := order.ApplyShipment(nil); err != nil {
		t.Fatal(err)
	}
	if order.Status() != OrderStatusShipped || order.Lines()[0].ShippedQuantity() != 2 {
		t.Fatalf("unexpected completed state: status=%s line=%+v", order.Status(), order.Lines()[0])
	}
}
