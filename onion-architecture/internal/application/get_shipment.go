package application

import "onion-architecture/internal/domain"

type ShipmentFinder interface {
	FindByID(id string) (domain.Shipment, error)
	ListByOrderID(orderID string) ([]domain.Shipment, error)
}

type GetShipmentQuery struct {
	ShipmentID string
}

type ShipmentDetails struct {
	ShipmentID string
	OrderID    string
	LineCount  int
}

type GetShipmentService struct {
	shipments ShipmentFinder
}

func NewGetShipmentService(shipments ShipmentFinder) GetShipmentService {
	return GetShipmentService{
		shipments: shipments,
	}
}

func (s GetShipmentService) Execute(query GetShipmentQuery) (ShipmentDetails, error) {
	shipment, err := s.shipments.FindByID(query.ShipmentID)
	if err != nil {
		return ShipmentDetails{}, err
	}

	return ShipmentDetails{
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
		LineCount:  len(shipment.Lines),
	}, nil
}
