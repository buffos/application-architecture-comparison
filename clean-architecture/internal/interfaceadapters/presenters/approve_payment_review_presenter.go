package presenters

import (
	"fmt"

	"clean-architecture/internal/usecases"
)

type ApprovePaymentReviewViewModel struct {
	Message string
	OrderID string
	Status  string
	Lines   int
}

type ApprovePaymentReviewPresenter struct {
	viewModel ApprovePaymentReviewViewModel
}

func NewApprovePaymentReviewPresenter() *ApprovePaymentReviewPresenter {
	return &ApprovePaymentReviewPresenter{}
}

func (p *ApprovePaymentReviewPresenter) Present(output usecases.ApprovePaymentReviewOutput) error {
	p.viewModel = ApprovePaymentReviewViewModel{
		Message: fmt.Sprintf("approved payment review: id=%s lines=%d status=%s", output.OrderID, output.Lines, output.Status),
		OrderID: output.OrderID,
		Status:  output.Status,
		Lines:   output.Lines,
	}

	return nil
}

func (p *ApprovePaymentReviewPresenter) ViewModel() ApprovePaymentReviewViewModel {
	return p.viewModel
}
