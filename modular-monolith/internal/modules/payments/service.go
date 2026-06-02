package payments

type Processor interface {
	Capture(request PaymentRequest) error
}

type Service struct {
	gateway Gateway
}

func NewService(gateway Gateway) Service {
	return Service{
		gateway: gateway,
	}
}

func (s Service) Capture(request PaymentRequest) error {
	return s.gateway.Capture(request)
}
