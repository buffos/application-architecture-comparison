package usecases

type ApprovePaymentReviewInput struct {
	OrderID string
}

type ApprovePaymentReviewOutput struct {
	OrderID string
	Status  string
	Lines   int
}

type ApprovePaymentReviewInputBoundary interface {
	Execute(input ApprovePaymentReviewInput) error
}

type ApprovePaymentReviewOutputBoundary interface {
	Present(output ApprovePaymentReviewOutput) error
}

type ApprovePaymentReviewInteractor struct {
	orders OrderEditor
	output ApprovePaymentReviewOutputBoundary
}

func NewApprovePaymentReviewInteractor(orders OrderEditor, output ApprovePaymentReviewOutputBoundary) ApprovePaymentReviewInteractor {
	return ApprovePaymentReviewInteractor{
		orders: orders,
		output: output,
	}
}

func (uc ApprovePaymentReviewInteractor) Execute(input ApprovePaymentReviewInput) error {
	order, err := uc.orders.FindByID(input.OrderID)
	if err != nil {
		return err
	}

	if err := order.ApprovePaymentReview(); err != nil {
		return err
	}

	if err := uc.orders.Save(order); err != nil {
		return err
	}

	return uc.output.Present(ApprovePaymentReviewOutput{
		OrderID: order.ID,
		Status:  order.Status,
		Lines:   len(order.Lines),
	})
}
