package shipments

type Creator interface {
	Create(request ShipmentRequest) (Shipment, error)
}

type Service struct {
	shipments Repository
}

func NewService(shipments Repository) Service {
	return Service{
		shipments: shipments,
	}
}

func (s Service) Create(request ShipmentRequest) (Shipment, error) {
	shipment := NewShipment(request)
	if err := s.shipments.Save(shipment); err != nil {
		return Shipment{}, err
	}

	return shipment, nil
}
