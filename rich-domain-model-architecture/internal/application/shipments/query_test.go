package shipments

import (
	"errors"
	"testing"

	"rich-domain-model-architecture/internal/domain/fulfillment"
	"rich-domain-model-architecture/internal/domain/ordering"
	"rich-domain-model-architecture/internal/domain/quoting"
)

func TestShipmentQueryProjectsDetailsAndFiltersSummaries(t *testing.T) {
	shipment := queryShipment(t)
	reader := NewInMemoryReader()
	reader.Save(shipment)
	details, err := reader.GetShipment(string(shipment.ID()))
	if err != nil {
		t.Fatal(err)
	}
	if details.Status != string(fulfillment.ShipmentStatusDispatched) || details.OrderID != "order-query" || len(details.Lines) != 1 {
		t.Fatalf("details = %+v", details)
	}
	rows := reader.ListShipments(fulfillment.ShipmentStatusDispatched)
	if len(rows) != 1 || rows[0].ShipmentID != string(shipment.ID()) {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestShipmentQueryReturnsNotFound(t *testing.T) {
	if _, err := NewInMemoryReader().GetShipment("missing"); !errors.Is(err, ErrShipmentNotFound) {
		t.Fatalf("missing query returned %v", err)
	}
}

func queryShipment(t *testing.T) fulfillment.Shipment {
	t.Helper()
	order := paidOrderForQuery(t)
	shipment, err := fulfillment.NewShipmentFromPaidOrder("shipment-query", order)
	if err != nil {
		t.Fatal(err)
	}
	if err := shipment.Dispatch(); err != nil {
		t.Fatal(err)
	}
	return shipment
}

func paidOrderForQuery(t *testing.T) ordering.Order {
	t.Helper()
	price, err := quoting.NewMoney(1000, "USD")
	if err != nil {
		t.Fatal(err)
	}
	line, err := quoting.NewQuoteLine("sku-001", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	quote, err := quoting.NewQuote("quote-shipment-query", "customer-001")
	if err != nil {
		t.Fatal(err)
	}
	if err := quote.AddLine(line); err != nil {
		t.Fatal(err)
	}
	if err := quote.SubmitForApproval(quoting.ApprovalDecision{}); err != nil {
		t.Fatal(err)
	}
	order, err := ordering.NewOrderFromQuote("order-query", quote)
	if err != nil {
		t.Fatal(err)
	}
	if err := order.MarkPaid(); err != nil {
		t.Fatal(err)
	}
	return order
}
