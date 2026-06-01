package application

type ListShipmentsQuery struct {
	OrderID string
}

type ListShipmentsService struct {
	shipments ShipmentFinder
}

func NewListShipmentsService(shipments ShipmentFinder) ListShipmentsService {
	return ListShipmentsService{
		shipments: shipments,
	}
}

func (s ListShipmentsService) Execute(query ListShipmentsQuery) ([]ShipmentDetails, error) {
	shipments, err := s.shipments.ListByOrderID(query.OrderID)
	if err != nil {
		return nil, err
	}

	result := make([]ShipmentDetails, 0, len(shipments))
	for _, shipment := range shipments {
		result = append(result, ShipmentDetails{
			ShipmentID: shipment.ID,
			OrderID:    shipment.OrderID,
			LineCount:  len(shipment.Lines),
		})
	}

	return result, nil
}
