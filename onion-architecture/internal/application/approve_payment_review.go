package application

type ApprovePaymentReviewCommand struct {
	OrderID string
}

type ApprovePaymentReviewResult struct {
	OrderID    string
	QuoteID    string
	CustomerID string
	Status     string
	LineCount  int
}

type ApprovePaymentReviewService struct {
	orders OrderRepository
}

func NewApprovePaymentReviewService(orders OrderRepository) ApprovePaymentReviewService {
	return ApprovePaymentReviewService{
		orders: orders,
	}
}

func (s ApprovePaymentReviewService) Execute(command ApprovePaymentReviewCommand) (ApprovePaymentReviewResult, error) {
	order, err := s.orders.FindByID(command.OrderID)
	if err != nil {
		return ApprovePaymentReviewResult{}, err
	}

	if err := order.ApprovePaymentReview(); err != nil {
		return ApprovePaymentReviewResult{}, err
	}

	if err := s.orders.Save(order); err != nil {
		return ApprovePaymentReviewResult{}, err
	}

	return ApprovePaymentReviewResult{
		OrderID:    order.ID,
		QuoteID:    order.QuoteID,
		CustomerID: order.CustomerID,
		Status:     order.Status,
		LineCount:  len(order.Lines),
	}, nil
}
