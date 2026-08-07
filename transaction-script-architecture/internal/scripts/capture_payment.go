package scripts

import (
	"errors"
	"fmt"
	"strings"

	"transaction-script-architecture/internal/data"
)

const (
	PaymentOutcomeAccept = "accept"
	PaymentOutcomeFail   = "fail"
	PaymentOutcomeReview = "review"
)

var (
	ErrOrderIDRequired       = errors.New("order id is required")
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderNotPayable       = errors.New("order is not ready for payment")
	ErrPaymentOutcomeInvalid = errors.New("payment outcome is invalid")
)

// CapturePayment records a simulated payment outcome and updates the order
// lifecycle in one procedural transaction.
func CapturePayment(store *data.Store, orderID string, outcome string) (data.Order, error) {
	if store == nil {
		return data.Order{}, ErrStoreRequired
	}

	if orderID == "" {
		return data.Order{}, ErrOrderIDRequired
	}

	order, ok := store.Orders[orderID]
	if !ok {
		return data.Order{}, ErrOrderNotFound
	}

	if order.Status != data.OrderStatusReadyForPayment {
		return data.Order{}, ErrOrderNotPayable
	}

	outcome = strings.ToLower(strings.TrimSpace(outcome))
	if outcome == "" {
		outcome = PaymentOutcomeAccept
	}

	status := data.PaymentStatusAccepted
	orderStatus := data.OrderStatusReadyForFulfillment
	switch outcome {
	case PaymentOutcomeAccept:
		status = data.PaymentStatusAccepted
	case PaymentOutcomeFail:
		status = data.PaymentStatusFailed
		orderStatus = data.OrderStatusReadyForPayment
	case PaymentOutcomeReview:
		status = data.PaymentStatusManualReview
		orderStatus = data.OrderStatusPaymentReview
	default:
		return data.Order{}, ErrPaymentOutcomeInvalid
	}

	store.NextPaymentNumber++
	payment := data.Payment{
		ID:      fmt.Sprintf("payment-%03d", store.NextPaymentNumber),
		OrderID: order.ID,
		Amount:  order.Total,
		Status:  status,
	}
	store.Payments[payment.ID] = payment

	order.PaymentID = payment.ID
	order.PaymentStatus = status
	order.Status = orderStatus
	store.Orders[order.ID] = order

	return order, nil
}
