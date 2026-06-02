package shipments

import "testing"

type stubRepository struct {
	saved Shipment
}

func (r *stubRepository) Save(shipment Shipment) error {
	r.saved = shipment
	return nil
}

func (r *stubRepository) FindByID(id string) (Shipment, error) {
	return r.saved, nil
}

func (r *stubRepository) ListByOrderID(orderID string) ([]Shipment, error) {
	if r.saved.ID == "" {
		return nil, nil
	}
	if orderID == "" || r.saved.OrderID == orderID {
		return []Shipment{r.saved}, nil
	}
	return nil, nil
}

func TestCreateSavesShipment(t *testing.T) {
	repository := &stubRepository{}
	service := NewService(repository)

	shipment, err := service.Create(ShipmentRequest{
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Lines: []ShipmentLine{
			{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if shipment.OrderID != "order-001" {
		t.Fatalf("expected order-001, got %s", shipment.OrderID)
	}

	if repository.saved.ID == "" {
		t.Fatalf("expected shipment to be saved")
	}
}
