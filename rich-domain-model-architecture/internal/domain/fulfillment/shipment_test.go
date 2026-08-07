package fulfillment

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/ordering"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func paidOrder(t *testing.T) ordering.Order {
	t.Helper()
	order, err := ordering.NewOrderFromQuote("order-001", approvedQuoteForShipment(t))
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	return order
}

func approvedQuoteForShipment(t *testing.T) quoting.Quote {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, price)
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
	return quote
}

func TestShipmentRequiresPaidOrderAndOwnsDispatchLifecycle(t *testing.T) {
	order := paidOrder(t)
	shipment, err := NewShipmentFromPaidOrder("shipment-001", order)
	if err != nil {
		t.Fatal(err)
	}
	if shipment.Status() != ShipmentStatusPending || len(shipment.Lines()) != 1 {
		t.Fatalf("unexpected shipment: status=%s lines=%d", shipment.Status(), len(shipment.Lines()))
	}
	if err := shipment.Dispatch(); err != nil {
		t.Fatal(err)
	}
	if shipment.Status() != ShipmentStatusDispatched {
		t.Fatalf("status = %s, want %s", shipment.Status(), ShipmentStatusDispatched)
	}
	if err := shipment.Dispatch(); !errors.Is(err, ErrShipmentNotDispatchable) {
		t.Fatalf("repeated dispatch returned %v", err)
	}
}

func TestShipmentRejectsUnpaidOrder(t *testing.T) {
	quote := approvedQuoteForShipment(t)
	order, err := ordering.NewOrderFromQuote("order-001", quote)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := NewShipmentFromPaidOrder("shipment-001", order); !errors.Is(err, ErrOrderNotPaid) {
		t.Fatalf("unpaid shipment creation returned %v", err)
	}
}

func TestShipmentCanContainPartialSelection(t *testing.T) {
	order := paidOrder(t)
	selection := []ordering.ShipmentSelection{{ProductSKU: "sku-001", Quantity: 1}}
	shipment, err := NewShipmentFromOrderSelection("shipment-001", order, selection)
	if err != nil {
		t.Fatal(err)
	}
	if len(shipment.Lines()) != 1 || shipment.Lines()[0].Quantity() != 1 {
		t.Fatalf("unexpected shipment: %+v", shipment)
	}
	if err := shipment.Dispatch(); err != nil {
		t.Fatal(err)
	}
	if err := order.ApplyShipment(selection); err != nil {
		t.Fatal(err)
	}
	if order.Status() != ordering.OrderStatusPartiallyShipped || order.Lines()[0].ShippedQuantity() != 1 {
		t.Fatalf("unexpected partial order: status=%s line=%+v", order.Status(), order.Lines()[0])
	}

	remaining, err := NewShipmentFromOrderSelection("shipment-002", order, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(remaining.Lines()) != 1 || remaining.Lines()[0].Quantity() != 1 {
		t.Fatalf("unexpected remaining shipment: %+v", remaining)
	}
	if err := order.ApplyShipment(nil); err != nil {
		t.Fatal(err)
	}
	if order.Status() != ordering.OrderStatusShipped || order.Lines()[0].ShippedQuantity() != 2 {
		t.Fatalf("unexpected completed order: status=%s line=%+v", order.Status(), order.Lines()[0])
	}
	if err := order.Cancel(); !errors.Is(err, ordering.ErrOrderNotCancellable) {
		t.Fatalf("partial shipment cancellation returned %v", err)
	}
}

func TestShipmentRejectsQuantityBeyondRemainingOrderQuantity(t *testing.T) {
	order := paidOrder(t)
	_, err := NewShipmentFromOrderSelection("shipment-001", order, []ordering.ShipmentSelection{{ProductSKU: "sku-001", Quantity: 3}})
	if !errors.Is(err, ErrShipmentSelectionInvalid) {
		t.Fatalf("invalid selection returned %v", err)
	}
}
