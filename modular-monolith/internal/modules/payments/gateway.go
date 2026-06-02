package payments

type Gateway interface {
	Capture(request PaymentRequest) (CaptureResult, error)
	Refund(request RefundRequest) error
}
