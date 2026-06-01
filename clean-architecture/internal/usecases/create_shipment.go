package usecases

import "clean-architecture/internal/entities"

type CreateShipmentInput struct {
	OrderID string
}

type CreateShipmentOutput struct {
	ShipmentID string
	OrderID    string
	Status     string
	Lines      int
}

type CreateShipmentInputBoundary interface {
	Execute(input CreateShipmentInput) error
}

type CreateShipmentOutputBoundary interface {
	Present(output CreateShipmentOutput) error
}

type ShipmentWriter interface {
	Save(shipment entities.Shipment) error
}

type CreateShipmentInteractor struct {
	orders    OrderEditor
	shipments ShipmentWriter
	clock     Clock
	output    CreateShipmentOutputBoundary
}

func NewCreateShipmentInteractor(orders OrderEditor, shipments ShipmentWriter, clock Clock, output CreateShipmentOutputBoundary) CreateShipmentInteractor {
	return CreateShipmentInteractor{
		orders:    orders,
		shipments: shipments,
		clock:     clock,
		output:    output,
	}
}

func (uc CreateShipmentInteractor) Execute(input CreateShipmentInput) error {
	order, err := uc.orders.FindByID(input.OrderID)
	if err != nil {
		return err
	}

	shipment, err := entities.NewShipmentFromPaidOrder(order)
	if err != nil {
		return err
	}

	if err := order.MarkShippedAt(uc.clock.Now()); err != nil {
		return err
	}

	if err := uc.shipments.Save(shipment); err != nil {
		return err
	}

	if err := uc.orders.Save(order); err != nil {
		return err
	}

	return uc.output.Present(CreateShipmentOutput{
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
		Status:     shipment.Status,
		Lines:      len(shipment.Lines),
	})
}
