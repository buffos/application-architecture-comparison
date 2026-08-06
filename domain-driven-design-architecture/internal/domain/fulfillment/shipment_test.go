package fulfillment

import (
	"errors"
	"testing"

	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
)

func paidOrder(t *testing.T) ordering.Order {
	t.Helper()
	quote, err := quoting.NewQuote("quote-001", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 2, price)
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
	return order
}

func TestShipmentRequiresPaidOrderAndDispatches(t *testing.T) {
	shipment, err := NewShipmentFromPaidOrder("shipment-001", paidOrder(t))
	if err != nil {
		t.Fatal(err)
	}
	if shipment.Status() != ShipmentStatusPending || len(shipment.Lines()) != 1 {
		t.Fatalf("unexpected shipment %+v", shipment)
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
	if _, err := NewShipmentFromPaidOrder("shipment-001", order); !errors.Is(err, ErrOrderNotPaid) {
		t.Fatalf("unpaid order returned %v", err)
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
		t.Fatalf("unexpected shipment %+v", shipment)
	}
}
