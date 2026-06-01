package usecases

import (
	"testing"

	"clean-architecture/internal/entities"
)

type stubShipmentReader struct {
	shipment entities.Shipment
	err      error
}

func (g stubShipmentReader) FindByID(id string) (entities.Shipment, error) {
	if g.err != nil {
		return entities.Shipment{}, g.err
	}

	return g.shipment, nil
}

type stubShipmentLister struct {
	shipments []entities.Shipment
	err       error
	orderID   string
}

func (g *stubShipmentLister) ListByOrderID(orderID string) ([]entities.Shipment, error) {
	g.orderID = orderID
	if g.err != nil {
		return nil, g.err
	}

	return g.shipments, nil
}

type stubGetShipmentOutput struct {
	output GetShipmentOutput
}

func (o *stubGetShipmentOutput) Present(output GetShipmentOutput) error {
	o.output = output
	return nil
}

type stubListShipmentsOutput struct {
	output ListShipmentsOutput
}

func (o *stubListShipmentsOutput) Present(output ListShipmentsOutput) error {
	o.output = output
	return nil
}

func TestGetShipmentInteractorLoadsShipment(t *testing.T) {
	output := &stubGetShipmentOutput{}
	interactor := NewGetShipmentInteractor(stubShipmentReader{
		shipment: entities.Shipment{
			ID:      "shipment-001",
			OrderID: "order-001",
			Status:  entities.ShipmentStatusCreated,
			Lines: []entities.ShipmentLine{
				{SKU: "CHAIR-001", Quantity: 2},
			},
		},
	}, output)

	err := interactor.Execute(GetShipmentInput{ShipmentID: "shipment-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if output.output.ShipmentID != "shipment-001" {
		t.Fatalf("expected shipment id shipment-001, got %s", output.output.ShipmentID)
	}

	if output.output.Lines != 1 {
		t.Fatalf("expected 1 line, got %d", output.output.Lines)
	}
}

func TestListShipmentsInteractorFiltersByOrderID(t *testing.T) {
	shipments := &stubShipmentLister{
		shipments: []entities.Shipment{
			{
				ID:      "shipment-001",
				OrderID: "order-001",
				Status:  entities.ShipmentStatusCreated,
			},
			{
				ID:      "shipment-002",
				OrderID: "order-001",
				Status:  entities.ShipmentStatusCreated,
			},
		},
	}
	output := &stubListShipmentsOutput{}
	interactor := NewListShipmentsInteractor(shipments, output)

	err := interactor.Execute(ListShipmentsInput{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if shipments.orderID != "order-001" {
		t.Fatalf("expected order filter order-001, got %s", shipments.orderID)
	}

	if output.output.Count != 2 {
		t.Fatalf("expected 2 shipments, got %d", output.output.Count)
	}

	if output.output.Shipments[0].ShipmentID != "shipment-001" {
		t.Fatalf("expected first shipment id shipment-001, got %s", output.output.Shipments[0].ShipmentID)
	}
}
