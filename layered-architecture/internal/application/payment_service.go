package application

import "layered-architecture/internal/domain"

type PaymentService struct {
	orderRepo OrderRepository
}

func NewPaymentService(orderRepo OrderRepository) PaymentService {
	return PaymentService{orderRepo: orderRepo}
}

func (s PaymentService) CapturePayment(orderID string) (domain.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return domain.Order{}, err
	}

	if err := order.AcceptPayment(); err != nil {
		return domain.Order{}, err
	}

	if err := s.orderRepo.Save(order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
