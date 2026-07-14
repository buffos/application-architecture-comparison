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
