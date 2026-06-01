package usecases

import "clean-architecture/internal/entities"

type ListShipmentsInput struct {
	OrderID string
}

type ShipmentListItem struct {
	ShipmentID string
	OrderID    string
	Status     string
	Lines      int
}

type ListShipmentsOutput struct {
	OrderID   string
	Count     int
	Shipments []ShipmentListItem
}

type ListShipmentsInputBoundary interface {
	Execute(input ListShipmentsInput) error
}

type ListShipmentsOutputBoundary interface {
	Present(output ListShipmentsOutput) error
}

type ShipmentLister interface {
	ListByOrderID(orderID string) ([]entities.Shipment, error)
}

type ListShipmentsInteractor struct {
	shipments ShipmentLister
	output    ListShipmentsOutputBoundary
}

func NewListShipmentsInteractor(shipments ShipmentLister, output ListShipmentsOutputBoundary) ListShipmentsInteractor {
	return ListShipmentsInteractor{
		shipments: shipments,
		output:    output,
	}
}

func (uc ListShipmentsInteractor) Execute(input ListShipmentsInput) error {
	shipments, err := uc.shipments.ListByOrderID(input.OrderID)
	if err != nil {
		return err
	}

	items := make([]ShipmentListItem, 0, len(shipments))
	for _, shipment := range shipments {
		items = append(items, ShipmentListItem{
			ShipmentID: shipment.ID,
			OrderID:    shipment.OrderID,
			Status:     shipment.Status,
			Lines:      len(shipment.Lines),
		})
	}

	return uc.output.Present(ListShipmentsOutput{
		OrderID:   input.OrderID,
		Count:     len(items),
		Shipments: items,
	})
}
