package scripts

import (
	"errors"
	"strings"

	"transaction-script-architecture/internal/data"
)

const (
	PaymentReviewDecisionAccept = "accept"
	PaymentReviewDecisionReject = "reject"
)

var (
	ErrOrderNotInPaymentReview = errors.New("order is not in payment review")
	ErrPaymentReviewMissing    = errors.New("payment review record not found")
	ErrPaymentReviewerRequired = errors.New("payment reviewer is required")
	ErrPaymentDecisionInvalid  = errors.New("payment review decision is invalid")
)

func ApprovePaymentReview(store *data.Store, orderID string, reviewedBy string, decision string, comment string) (data.Order, error) {
	if store == nil {
		return data.Order{}, ErrStoreRequired
	}
	if orderID == "" {
		return data.Order{}, ErrOrderIDRequired
	}
	if reviewedBy == "" {
		return data.Order{}, ErrPaymentReviewerRequired
	}

	order, ok := store.Orders[orderID]
	if !ok {
		return data.Order{}, ErrOrderNotFound
	}
	if order.Status != data.OrderStatusPaymentReview {
		return data.Order{}, ErrOrderNotInPaymentReview
	}

	payment, ok := store.Payments[order.PaymentID]
	if !ok || payment.Status != data.PaymentStatusManualReview {
		return data.Order{}, ErrPaymentReviewMissing
	}

	switch strings.ToLower(strings.TrimSpace(decision)) {
	case PaymentReviewDecisionAccept:
		payment.Status = data.PaymentStatusAccepted
		order.Status = data.OrderStatusReadyForFulfillment
		order.PaymentStatus = data.PaymentStatusAccepted
	case PaymentReviewDecisionReject:
		payment.Status = data.PaymentStatusFailed
		order.Status = data.OrderStatusReadyForPayment
		order.PaymentStatus = data.PaymentStatusFailed
	default:
		return data.Order{}, ErrPaymentDecisionInvalid
	}

	payment.ReviewedBy = reviewedBy
	payment.DecisionComment = comment
	store.Payments[payment.ID] = payment
	store.Orders[order.ID] = order

	return order, nil
}
