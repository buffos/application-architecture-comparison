package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type CapturePaymentUseCase struct {
	orders   ports.OrderRepository
	payments ports.PaymentGateway
}

func NewCapturePaymentUseCase(orders ports.OrderRepository, payments ports.PaymentGateway) CapturePaymentUseCase {
	return CapturePaymentUseCase{
		orders:   orders,
		payments: payments,
	}
}

func (uc CapturePaymentUseCase) Execute(id string) (domain.Order, error) {
	order, err := uc.orders.FindByID(id)
	if err != nil {
		return domain.Order{}, err
	}

	ok, err := uc.payments.Capture(order)
	if err != nil {
		return domain.Order{}, err
	}

	if !ok {
		return domain.Order{}, domain.ErrPaymentFailed
	}

	order.AcceptPayment()

	if err := uc.orders.Save(order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
