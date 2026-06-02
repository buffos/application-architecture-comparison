package shipments

import "testing"

type stubQueryRepository struct {
	shipments map[string]Shipment
}

func (r *stubQueryRepository) Save(shipment Shipment) error {
	if r.shipments == nil {
		r.shipments = make(map[string]Shipment)
	}
	r.shipments[shipment.ID] = shipment
	return nil
}

func (r *stubQueryRepository) FindByID(id string) (Shipment, error) {
	shipment, ok := r.shipments[id]
	if !ok {
		return Shipment{}, ErrShipmentNotFound
	}
	return shipment, nil
}

func (r *stubQueryRepository) ListByOrderID(orderID string) ([]Shipment, error) {
	list := make([]Shipment, 0, len(r.shipments))
	for _, shipment := range r.shipments {
		if orderID == "" || shipment.OrderID == orderID {
			list = append(list, shipment)
		}
	}
	return list, nil
}

func TestGetShipmentLoadsStoredShipment(t *testing.T) {
	repository := &stubQueryRepository{shipments: map[string]Shipment{}}
	shipment := Shipment{
		ID:         "shipment-001",
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Lines: []ShipmentLine{
			{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 2},
		},
	}
	_ = repository.Save(shipment)
	service := NewService(repository)

	result, err := service.GetShipment(GetShipmentQuery{ShipmentID: shipment.ID})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ShipmentID != shipment.ID || result.OrderID != shipment.OrderID {
		t.Fatalf("expected stored shipment details to be returned")
	}
}

func TestListShipmentsFiltersByOrderID(t *testing.T) {
	repository := &stubQueryRepository{shipments: map[string]Shipment{}}
	_ = repository.Save(Shipment{
		ID:         "shipment-001",
		OrderID:    "order-001",
		CustomerID: "customer-001",
		Lines:      []ShipmentLine{{ProductSKU: "sku-001", ProductName: "Desk", Quantity: 2}},
	})
	_ = repository.Save(Shipment{
		ID:         "shipment-002",
		OrderID:    "order-002",
		CustomerID: "customer-002",
		Lines:      []ShipmentLine{{ProductSKU: "sku-002", ProductName: "Chair", Quantity: 1}},
	})
	service := NewService(repository)

	result, err := service.ListShipments(ListShipmentsQuery{OrderID: "order-001"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result) != 1 || result[0].OrderID != "order-001" {
		t.Fatalf("expected one shipment for order-001, got %+v", result)
	}
}
