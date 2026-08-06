package shipments

import (
	"testing"

	"domain-driven-design-architecture/internal/domain/fulfillment"
	"domain-driven-design-architecture/internal/domain/ordering"
	"domain-driven-design-architecture/internal/domain/quoting"
)

func TestReaderProjectsShipmentAggregate(t *testing.T) {
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
	shipment, err := fulfillment.NewShipmentFromPaidOrder("shipment-001", order)
	if err != nil {
		t.Fatal(err)
	}
	if err := shipment.Dispatch(); err != nil {
		t.Fatal(err)
	}
	reader := NewInMemoryReader()
	reader.Save(shipment)
	details, err := reader.GetShipment("shipment-001")
	if err != nil {
		t.Fatal(err)
	}
	if details.Status != string(fulfillment.ShipmentStatusDispatched) || len(details.Lines) != 1 || details.Lines[0].Quantity != 2 {
		t.Fatalf("unexpected details %+v", details)
	}
	if got := reader.ListShipments(fulfillment.ShipmentStatusDispatched); len(got) != 1 {
		t.Fatalf("unexpected summaries %+v", got)
	}
}
