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

	result, err := uc.payments.Capture(order)
	if err != nil {
		return domain.Order{}, err
	}

	switch result {
	case ports.PaymentResultAccepted:
		order.AcceptPayment()
	case ports.PaymentResultManualReview:
		order.MarkPaymentReview()
	case ports.PaymentResultFailed:
		order.FailPayment()
		if err := uc.orders.Save(order); err != nil {
			return domain.Order{}, err
		}
		return domain.Order{}, domain.ErrPaymentFailed
	default:
		order.FailPayment()
		if err := uc.orders.Save(order); err != nil {
			return domain.Order{}, err
		}
		return domain.Order{}, domain.ErrPaymentFailed
	}

	if err := uc.orders.Save(order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
