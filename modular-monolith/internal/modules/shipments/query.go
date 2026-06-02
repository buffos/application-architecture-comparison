package shipments

type GetShipmentQuery struct {
	ShipmentID string
}

type ShipmentDetails struct {
	ShipmentID string
	OrderID    string
	CustomerID string
	LineCount  int
}

type ListShipmentsQuery struct {
	OrderID string
}

func (s Service) GetShipment(query GetShipmentQuery) (ShipmentDetails, error) {
	shipment, err := s.shipments.FindByID(query.ShipmentID)
	if err != nil {
		return ShipmentDetails{}, err
	}

	return ShipmentDetails{
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
		CustomerID: shipment.CustomerID,
		LineCount:  len(shipment.Lines),
	}, nil
}

func (s Service) ListShipments(query ListShipmentsQuery) ([]ShipmentDetails, error) {
	shipments, err := s.shipments.ListByOrderID(query.OrderID)
	if err != nil {
		return nil, err
	}

	list := make([]ShipmentDetails, 0, len(shipments))
	for _, shipment := range shipments {
		list = append(list, ShipmentDetails{
			ShipmentID: shipment.ID,
			OrderID:    shipment.OrderID,
			CustomerID: shipment.CustomerID,
			LineCount:  len(shipment.Lines),
		})
	}

	return list, nil
}
