package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"hexagonal-architecture/internal/core/application"
	"hexagonal-architecture/internal/core/domain"
)

type PaymentHandler struct {
	capturePayment       application.CapturePaymentUseCase
	approvePaymentReview application.ApprovePaymentReviewUseCase
}

func NewPaymentHandler(
	capturePayment application.CapturePaymentUseCase,
	approvePaymentReview application.ApprovePaymentReviewUseCase,
) PaymentHandler {
	return PaymentHandler{
		capturePayment:       capturePayment,
		approvePaymentReview: approvePaymentReview,
	}
}

func (h PaymentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/capture-payment"):
		h.capture(w, r)
	case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/approve-payment-review"):
		h.approveReview(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h PaymentHandler) capture(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/orders/"), "/capture-payment")

	order, err := h.capturePayment.Execute(orderID)
	if err != nil {
		http.Error(w, err.Error(), paymentStatusCode(err))
		return
	}

	writeOrderResponse(w, http.StatusOK, order)
}

func (h PaymentHandler) approveReview(w http.ResponseWriter, r *http.Request) {
	orderID := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/orders/"), "/approve-payment-review")

	var body struct {
		ReviewedBy string `json:"reviewedBy"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := h.approvePaymentReview.Execute(orderID, body.ReviewedBy)
	if err != nil {
		http.Error(w, err.Error(), paymentStatusCode(err))
		return
	}

	writeOrderResponse(w, http.StatusOK, order)
}

func writeOrderResponse(w http.ResponseWriter, status int, order domain.Order) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(toOrderResponse(order))
}

func paymentStatusCode(err error) int {
	switch err {
	case domain.ErrOrderNotFound:
		return http.StatusNotFound
	case domain.ErrPaymentFailed, domain.ErrPaymentReviewNotAllowed, domain.ErrOrderActorRequired:
		return http.StatusBadRequest
	default:
		return http.StatusBadRequest
	}
}
