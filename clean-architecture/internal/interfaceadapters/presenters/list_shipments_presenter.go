package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type ShipmentListItemViewModel struct {
	ShipmentID string
	OrderID    string
	Status     string
	Lines      int
}

type ListShipmentsViewModel struct {
	Message   string
	OrderID   string
	Count     int
	Shipments []ShipmentListItemViewModel
}

type ListShipmentsPresenter struct {
	viewModel ListShipmentsViewModel
}

func NewListShipmentsPresenter() *ListShipmentsPresenter {
	return &ListShipmentsPresenter{}
}

func (p *ListShipmentsPresenter) Present(output usecases.ListShipmentsOutput) error {
	items := make([]ShipmentListItemViewModel, 0, len(output.Shipments))
	for _, shipment := range output.Shipments {
		items = append(items, ShipmentListItemViewModel{
			ShipmentID: shipment.ShipmentID,
			OrderID:    shipment.OrderID,
			Status:     shipment.Status,
			Lines:      shipment.Lines,
		})
	}

	p.viewModel = ListShipmentsViewModel{
		Message:   fmt.Sprintf("listed shipments: order=%s count=%d", output.OrderID, output.Count),
		OrderID:   output.OrderID,
		Count:     output.Count,
		Shipments: items,
	}

	return nil
}

func (p *ListShipmentsPresenter) ViewModel() ListShipmentsViewModel {
	return p.viewModel
}
