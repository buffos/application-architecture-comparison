package shipments

import "microkernel-architecture/internal/kernel"

type Service struct {
	shipments Repository
}

func NewService(shipments Repository) Service {
	return Service{
		shipments: shipments,
	}
}

func (s Service) CreateShipment(record kernel.CreateShipmentRecord) (kernel.ShipmentCreationResult, error) {
	lines := make([]ShipmentLine, 0, len(record.Lines))
	for _, line := range record.Lines {
		lines = append(lines, ShipmentLine{
			ProductSKU: line.ProductSKU,
			Quantity:   line.Quantity,
		})
	}

	shipment := NewShipment(record.OrderID, record.CustomerID, lines)
	if err := s.shipments.Save(shipment); err != nil {
		return kernel.ShipmentCreationResult{}, err
	}

	return kernel.ShipmentCreationResult{
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
		CustomerID: shipment.CustomerID,
		LineCount:  len(shipment.Lines),
	}, nil
}

func (s Service) GetShipment(query kernel.GetShipmentQuery) (kernel.ShipmentDetails, error) {
	shipment, err := s.shipments.FindByID(query.ShipmentID)
	if err != nil {
		return kernel.ShipmentDetails{}, err
	}

	return kernel.ShipmentDetails{
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
		CustomerID: shipment.CustomerID,
		LineCount:  len(shipment.Lines),
	}, nil
}

func (s Service) ListShipments(query kernel.ListShipmentsQuery) ([]kernel.ShipmentSummary, error) {
	shipmentsList, err := s.shipments.ListByOrderID(query.OrderID)
	if err != nil {
		return nil, err
	}

	results := make([]kernel.ShipmentSummary, 0, len(shipmentsList))
	for _, shipment := range shipmentsList {
		results = append(results, kernel.ShipmentSummary{
			ShipmentID: shipment.ID,
			OrderID:    shipment.OrderID,
			CustomerID: shipment.CustomerID,
			LineCount:  len(shipment.Lines),
		})
	}

	return results, nil
}
