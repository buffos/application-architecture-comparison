package application

import "onion-architecture/internal/domain"

const PaymentCaptureOutcomeApproved = "Approved"
const PaymentCaptureOutcomeReview = "Review"

type CapturePaymentCommand struct {
	OrderID string
}

type CapturePaymentResult struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type OrderRepository interface {
	FindByID(id string) (domain.Order, error)
	Save(order domain.Order) error
}

type PaymentGateway interface {
	Capture(order domain.Order) (string, error)
}

type CapturePaymentService struct {
	orders   OrderRepository
	payments PaymentGateway
}

func NewCapturePaymentService(orders OrderRepository, payments PaymentGateway) CapturePaymentService {
	return CapturePaymentService{
		orders:   orders,
		payments: payments,
	}
}

func (s CapturePaymentService) Execute(command CapturePaymentCommand) (CapturePaymentResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return CapturePaymentResult{}, err
	}

	outcome, err := s.payments.Capture(order)
	if err != nil {
		return CapturePaymentResult{}, err
	}

	switch outcome {
	case PaymentCaptureOutcomeApproved:
		if err := order.MarkPaid(); err != nil {
			return CapturePaymentResult{}, err
		}
	case PaymentCaptureOutcomeReview:
		if err := order.MarkPaymentReview(); err != nil {
			return CapturePaymentResult{}, err
		}
	default:
		return CapturePaymentResult{}, domain.ErrOrderNotPayable
	}

	if err := s.orders.Save(order); err != nil {
		return CapturePaymentResult{}, err
	}

	return CapturePaymentResult{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}
