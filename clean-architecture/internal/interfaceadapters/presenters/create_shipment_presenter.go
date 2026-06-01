package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type CreateShipmentViewModel struct {
	Message    string
	ShipmentID string
	OrderID    string
	Status     string
	Lines      int
}

type CreateShipmentPresenter struct {
	viewModel CreateShipmentViewModel
}

func NewCreateShipmentPresenter() *CreateShipmentPresenter {
	return &CreateShipmentPresenter{}
}

func (p *CreateShipmentPresenter) Present(output usecases.CreateShipmentOutput) error {
	p.viewModel = CreateShipmentViewModel{
		Message:    fmt.Sprintf("created shipment: id=%s order=%s lines=%d status=%s", output.ShipmentID, output.OrderID, output.Lines, output.Status),
		ShipmentID: output.ShipmentID,
		OrderID:    output.OrderID,
		Status:     output.Status,
		Lines:      output.Lines,
	}

	return nil
}

func (p *CreateShipmentPresenter) ViewModel() CreateShipmentViewModel {
	return p.viewModel
}
