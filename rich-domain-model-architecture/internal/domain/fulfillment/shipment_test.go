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
	line, err := quoting.NewQuoteLine("sku-001", 1, price)
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
