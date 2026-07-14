package shipments

import (
	"testing"

	"microkernel-architecture/internal/kernel"
)

type stubRepository struct {
	saved Shipment
}

func (r *stubRepository) FindByID(id string) (Shipment, error) {
	if r.saved.ID == id {
		return r.saved, nil
	}

	return Shipment{}, ErrShipmentNotFound
}

func (r *stubRepository) ListByOrderID(orderID string) ([]Shipment, error) {
	if r.saved.ID == "" {
		return []Shipment{}, nil
	}

	if orderID == "" || r.saved.OrderID == orderID {
		return []Shipment{r.saved}, nil
	}

	return []Shipment{}, nil
}

func (r *stubRepository) Save(shipment Shipment) error {
	r.saved = shipment
	return nil
}

func TestCreateShipment(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository)

	result, err := service.CreateShipment(kernel.CreateShipmentRecord{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Lines: []kernel.ShipmentLine{
			{ProductSKU: "sku-002", Quantity: 1},
		},
	})
	if err != nil {
		t.Fatalf("expected create shipment to succeed, got %v", err)
	}

	if result.OrderID != "order-001" {
		t.Fatalf("expected order id order-001, got %s", result.OrderID)
	}
}

func TestGetShipment(t *testing.T) {
	repository := &stubRepository{
		saved: Shipment{
			ID:         "shipment-001",
			OrderID:    "order-001",
			CustomerID: "customer-001",
			Lines: []ShipmentLine{
				{ProductSKU: "sku-002", Quantity: 1},
			},
		},
	}
	service := NewService(repository)

	result, err := service.GetShipment(kernel.GetShipmentQuery{
		ShipmentID: "shipment-001",
	})
	if err != nil {
		t.Fatalf("expected get shipment to succeed, got %v", err)
	}

	if result.ShipmentID != "shipment-001" || result.OrderID != "order-001" {
		t.Fatalf("unexpected shipment details %+v", result)
	}
}

func TestListShipmentsByOrderID(t *testing.T) {
	repository := &stubRepository{
		saved: Shipment{
			ID:         "shipment-001",
			OrderID:    "order-001",
			CustomerID: "customer-001",
			Lines: []ShipmentLine{
				{ProductSKU: "sku-002", Quantity: 1},
			},
		},
	}
	service := NewService(repository)

	result, err := service.ListShipments(kernel.ListShipmentsQuery{
		OrderID: "order-001",
	})
	if err != nil {
		t.Fatalf("expected list shipments to succeed, got %v", err)
	}

	if len(result) != 1 || result[0].ShipmentID != "shipment-001" {
		t.Fatalf("unexpected shipment list %+v", result)
	}
}
