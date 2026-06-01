package usecases

import "clean-architecture/internal/entities"

type CapturePaymentInput struct {
	OrderID string
}

type CapturePaymentOutput struct {
	OrderID string
	Status  string
	Lines   int
}

type CapturePaymentInputBoundary interface {
	Execute(input CapturePaymentInput) error
}

type CapturePaymentOutputBoundary interface {
	Present(output CapturePaymentOutput) error
}

type OrderEditor interface {
	FindByID(id string) (entities.Order, error)
	Save(order entities.Order) error
}

type PaymentGateway interface {
	Capture(order entities.Order) error
}

type CapturePaymentInteractor struct {
	orders  OrderEditor
	payment PaymentGateway
	output  CapturePaymentOutputBoundary
}

func NewCapturePaymentInteractor(orders OrderEditor, payment PaymentGateway, output CapturePaymentOutputBoundary) CapturePaymentInteractor {
	return CapturePaymentInteractor{
		orders:  orders,
		payment: payment,
		output:  output,
	}
}

func (uc CapturePaymentInteractor) Execute(input CapturePaymentInput) error {
	order, err := uc.orders.FindByID(input.OrderID)
	if err != nil {
		return err
	}

	if err := uc.payment.Capture(order); err != nil {
		return err
	}

	if err := order.MarkPaid(); err != nil {
		return err
	}

	if err := uc.orders.Save(order); err != nil {
		return err
	}

	return uc.output.Present(CapturePaymentOutput{
		OrderID: order.ID,
		Status:  order.Status,
		Lines:   len(order.Lines),
	})
}
