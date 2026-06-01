package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type GetShipmentViewModel struct {
	Message    string
	ShipmentID string
	OrderID    string
	Status     string
	Lines      int
}

type GetShipmentPresenter struct {
	viewModel GetShipmentViewModel
}

func NewGetShipmentPresenter() *GetShipmentPresenter {
	return &GetShipmentPresenter{}
}

func (p *GetShipmentPresenter) Present(output usecases.GetShipmentOutput) error {
	p.viewModel = GetShipmentViewModel{
		Message:    fmt.Sprintf("loaded shipment: id=%s order=%s lines=%d status=%s", output.ShipmentID, output.OrderID, output.Lines, output.Status),
		ShipmentID: output.ShipmentID,
		OrderID:    output.OrderID,
		Status:     output.Status,
		Lines:      output.Lines,
	}

	return nil
}

func (p *GetShipmentPresenter) ViewModel() GetShipmentViewModel {
	return p.viewModel
}
