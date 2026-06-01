package application

import "onion-architecture/internal/domain"

type CreateShipmentCommand struct {
	OrderID string
}

type CreateShipmentResult struct {
	ShipmentID string
	OrderID    string
	OrderStatus string
	LineCount  int
}

type ShipmentStore interface {
	Save(shipment domain.Shipment) error
}

type CreateShipmentService struct {
	orders    OrderRepository
	shipments ShipmentStore
	clock     Clock
}

func NewCreateShipmentService(orders OrderRepository, shipments ShipmentStore, clock Clock) CreateShipmentService {
	return CreateShipmentService{
		orders:    orders,
		shipments: shipments,
		clock:     clock,
	}
}

func (s CreateShipmentService) Execute(command CreateShipmentCommand) (CreateShipmentResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return CreateShipmentResult{}, err
	}

	shipment, err := domain.NewShipmentFromOrder(order)
	if err != nil {
		return CreateShipmentResult{}, err
	}

	if err := order.MarkShipped(s.clock.Now()); err != nil {
		return CreateShipmentResult{}, err
	}

	if err := s.shipments.Save(shipment); err != nil {
		return CreateShipmentResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return CreateShipmentResult{}, err
	}

	return CreateShipmentResult{
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
		OrderStatus: order.Status,
		LineCount:  len(shipment.Lines),
	}, nil
}
