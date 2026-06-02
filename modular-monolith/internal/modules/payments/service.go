package payments

type Processor interface {
	Capture(request PaymentRequest) (CaptureResult, error)
}

type Refunder interface {
	Refund(request RefundRequest) error
}

type Service struct {
	gateway Gateway
}

func NewService(gateway Gateway) Service {
	return Service{
		gateway: gateway,
	}
}

func (s Service) Capture(request PaymentRequest) (CaptureResult, error) {
	return s.gateway.Capture(request)
}

func (s Service) Refund(request RefundRequest) error {
	return s.gateway.Refund(request)
}
