package application

import (
	"hexagonal-architecture/internal/core/domain"
	"hexagonal-architecture/internal/core/ports"
)

type ApprovePaymentReviewUseCase struct {
	orders ports.OrderRepository
}

func NewApprovePaymentReviewUseCase(orders ports.OrderRepository) ApprovePaymentReviewUseCase {
	return ApprovePaymentReviewUseCase{orders: orders}
}

func (uc ApprovePaymentReviewUseCase) Execute(orderID string, reviewedBy string) (domain.Order, error) {
	order, err := uc.orders.FindByID(orderID)
	if err != nil {
		return domain.Order{}, err
	}

	if err := order.ApprovePaymentReview(reviewedBy); err != nil {
		return domain.Order{}, err
	}

	if err := uc.orders.Save(order); err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
