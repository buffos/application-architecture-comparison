package payments

type Gateway interface {
	Capture(request PaymentRequest) error
	Refund(request RefundRequest) error
}
