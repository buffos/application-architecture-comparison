package payments

// Processor is the public contract this component provides to order workflows.
type Processor interface {
	Capture(request PaymentRequest) (CaptureResult, error)
}

type Refunder interface {
	Refund(request RefundRequest) error
}

// Gateway is the component's internal integration contract. Concrete adapters
// are selected at the composition root.
type Gateway interface {
	Capture(request PaymentRequest) (CaptureResult, error)
	Refund(request RefundRequest) error
}
