package application

import (
	"testing"

	"onion-architecture/internal/domain"
)

type stubShipmentFinder struct {
	shipment domain.Shipment
	list     []domain.Shipment
	err      error
}

func (f stubShipmentFinder) FindByID(id string) (domain.Shipment, error) {
	if f.err != nil {
		return domain.Shipment{}, f.err
	}

	return f.shipment, nil
}

func (f stubShipmentFinder) ListByOrderID(orderID string) ([]domain.Shipment, error) {
	if f.err != nil {
		return nil, f.err
	}

	result := make([]domain.Shipment, 0)
	for _, shipment := range f.list {
		if shipment.OrderID == orderID {
			result = append(result, shipment)
		}
	}

	return result, nil
}

func TestGetShipmentServiceReturnsShipmentDetails(t *testing.T) {
	service := NewGetShipmentService(stubShipmentFinder{
		shipment: domain.Shipment{
			ID:      "shipment-001",
			OrderID: "order-001",
			Lines: []domain.ShipmentLine{
				{ProductSKU: "sku-002", Quantity: 1},
			},
		},
	})

	result, err := service.Execute(GetShipmentQuery{ShipmentID: "shipment-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ShipmentID != "shipment-001" {
		t.Fatalf("expected shipment-001, got %s", result.ShipmentID)
	}
}

func TestListShipmentsServiceFiltersByOrderID(t *testing.T) {
	service := NewListShipmentsService(stubShipmentFinder{
		list: []domain.Shipment{
			{ID: "shipment-001", OrderID: "order-001"},
			{ID: "shipment-002", OrderID: "order-002"},
		},
	})

	result, err := service.Execute(ListShipmentsQuery{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result))
	}

	if result[0].ShipmentID != "shipment-001" {
		t.Fatalf("expected shipment-001, got %s", result[0].ShipmentID)
	}
}
