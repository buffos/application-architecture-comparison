package usecases

import "clean-architecture/internal/entities"

type CreateShipmentInput struct {
	OrderID string
	Lines   []CreateShipmentLineInput
}

type CreateShipmentLineInput struct {
	SKU      string
	Quantity int
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

	requestedLines := make([]entities.ShipmentLine, 0, len(input.Lines))
	for _, line := range input.Lines {
		productName := ""
		for _, orderLine := range order.Lines {
			if orderLine.SKU == line.SKU {
				productName = orderLine.ProductName
				break
			}
		}

		requestedLines = append(requestedLines, entities.ShipmentLine{
			SKU:         line.SKU,
			ProductName: productName,
			Quantity:    line.Quantity,
		})
	}

	shipment, err := entities.NewShipmentFromOrder(order, requestedLines)
	if err != nil {
		return err
	}

	if err := order.ApplyShipment(shipment.Lines, uc.clock.Now()); err != nil {
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
