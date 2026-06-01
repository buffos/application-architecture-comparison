package usecases

import "clean-architecture/internal/entities"

type GetShipmentInput struct {
	ShipmentID string
}

type GetShipmentOutput struct {
	ShipmentID string
	OrderID    string
	Status     string
	Lines      int
}

type GetShipmentInputBoundary interface {
	Execute(input GetShipmentInput) error
}

type GetShipmentOutputBoundary interface {
	Present(output GetShipmentOutput) error
}

type ShipmentReader interface {
	FindByID(id string) (entities.Shipment, error)
}

type GetShipmentInteractor struct {
	shipments ShipmentReader
	output    GetShipmentOutputBoundary
}

func NewGetShipmentInteractor(shipments ShipmentReader, output GetShipmentOutputBoundary) GetShipmentInteractor {
	return GetShipmentInteractor{
		shipments: shipments,
		output:    output,
	}
}

func (uc GetShipmentInteractor) Execute(input GetShipmentInput) error {
	shipment, err := uc.shipments.FindByID(input.ShipmentID)
	if err != nil {
		return err
	}

	return uc.output.Present(GetShipmentOutput{
		ShipmentID: shipment.ID,
		OrderID:    shipment.OrderID,
		Status:     shipment.Status,
		Lines:      len(shipment.Lines),
	})
}
